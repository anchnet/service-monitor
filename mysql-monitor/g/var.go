package g

import (
	"sync"
	"time"

	log "github.com/cihub/seelog"

	"math/rand"

	"github.com/open-falcon/common/model"
)

var (
	TransferClientsLock *sync.RWMutex                   = new(sync.RWMutex)
	TransferClients     map[string]*SingleConnRpcClient = map[string]*SingleConnRpcClient{}
)

func initTransferClient(addr string) {
	TransferClientsLock.Lock()
	defer TransferClientsLock.Unlock()
	TransferClients[addr] = &SingleConnRpcClient{
		RpcServer: addr,
		Timeout:   time.Duration(Config().Transfer.Timeout) * time.Millisecond,
	}
}

func SendToTransfer(metrics []*model.MetricValue) {
	if len(metrics) == 0 {
		return
	}

	debug := Config().Debug

	if debug {
		for i, _ := range metrics {
			log.Infof("=> <Total=%d> %v\n", len(metrics), metrics[i])
		}
	}

	var resp model.TransferResponse
	SendMetrics(metrics, &resp)
	if debug {
		log.Info("<=", &resp)
	}
}

func SendMetrics(metrics []*model.MetricValue, resp *model.TransferResponse) {
	rand.Seed(time.Now().UnixNano())
	for _, i := range rand.Perm(len(Config().Transfer.Addrs)) {
		addr := Config().Transfer.Addrs[i]
		if _, ok := TransferClients[addr]; !ok {
			initTransferClient(addr)
		}
		if updateMetrics(addr, metrics, resp) {
			break
		}
	}
}

func updateMetrics(addr string, metrics []*model.MetricValue, resp *model.TransferResponse) bool {
	TransferClientsLock.RLock()
	defer TransferClientsLock.RUnlock()
	err := TransferClients[addr].Call("Transfer.Update", metrics, resp)
	if err != nil {
		log.Info("call Transfer.Update fail", addr, err)
		return false
	}
	return true
}
