package funcs

import (
	log "github.com/cihub/seelog"

	"github.com/open-falcon/common/model"
	"github.com/toolkits/nux"
)

func UdpMetrics() []*model.MetricValue {
	udp, err := nux.Snmp("Udp")
	if err != nil {
		log.Info("read snmp fail", err)
		return []*model.MetricValue{}
	}

	count := len(udp)
	ret := make([]*model.MetricValue, count)
	i := 0
	for key, val := range udp {
		ret[i] = CounterValue("snmp.Udp."+key, val)
		i++
	}

	return ret
}
