package funcs

import (
	log "github.com/cihub/seelog"

	"github.com/open-falcon/common/model"
	"github.com/toolkits/nux"
)

func SocketStatSummaryMetrics() (L []*model.MetricValue) {
	ssMap, err := nux.SocketStatSummary()
	if err != nil {
		log.Info(err)
		return
	}

	for k, v := range ssMap {
		L = append(L, GaugeValue("ss."+k, v))
	}

	return
}
