package g

import (
	"log"
	"os"
	"time"

	"github.com/open-falcon/common/model"
)

var Root string

func InitRootDir() {
	var err error
	Root, err = os.Getwd()
	if err != nil {
		log.Fatalln("getwd fail:", err)
	}
}

var (
	TransferClient *SingleConnRpcClient
)

func InitRpcClients() {
	if Config().Transfer.Enabled {
		TransferClient = &SingleConnRpcClient{
			RpcServer: Config().Transfer.Addr,
			Timeout:   time.Duration(Config().Transfer.Timeout) * time.Millisecond,
		}
	}
}

func SendToTransfer(metrics []*model.MetricValue) {
	if len(metrics) == 0 {
		return
	}

	debug := Config().Debug

	if debug {
		for i, _ := range metrics {
			log.Printf("=> <Total=%d> %v\n", len(metrics), metrics[i])
		}
	}

	var resp model.TransferResponse
	err := TransferClient.Call("Transfer.Update", metrics, &resp)
	if err != nil {
		log.Println("call Transfer.Update fail", err)
	}

	if debug {
		log.Println("<=", &resp)
	}
}
