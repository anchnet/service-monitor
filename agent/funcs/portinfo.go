package funcs

import (
	"github.com/open-falcon/common/model"
	"log"
	"net"
	"time"
	"github.com/51idc/service-monitor/agent/g"
)

func PortMetricsCus() (L []*model.MetricValue) {
	port_list := g.Config().Port
	for _, port := range port_list {
		tag := "port=" + port
		L = append(L, GaugeValue("port.alive", port_dail(port), tag))
	}
	return L
}

func port_dail(tcp_port string) int {
	tcpaddr := "localhost:" + tcp_port
	// 40s time out
	conn, err := net.DialTimeout("tcp", tcpaddr, g.Config().DialTimeout * time.Second)
	if err != nil {
		log.Println("Port:" + tcp_port + ".DialTimeout error :", err)
		return 0
	}
	conn.Close()
	return 1
}