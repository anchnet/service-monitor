package funcs

import (
	log "github.com/cihub/seelog"
	"os/exec"
	"strings"

	"github.com/open-falcon/common/model"
)

func GetNetStat() (map[string]uint64, error) {
	var timewait uint64 = 0
	var closed uint64 = 0
	var listen uint64 = 0
	var syn uint64 = 0
	var estab uint64 = 0
	m := make(map[string]uint64)

	out, err := exec.Command("netstat", "-an").Output()
	if err != nil {
		return nil, err
	}
	netstat := strings.Split(string(out), "\r\n")
	for _, line := range netstat {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "TCP") {
			if strings.Contains(l, "TIME_WAIT") {
				timewait += 1
			}
			if strings.Contains(l, "ESTABLISHED") {
				estab += 1
			}
			if strings.Contains(l, "CLOSE_WAIT") {
				closed += 1
			}
			if strings.Contains(l, "SYN") || strings.Contains(l, "FIN") || strings.Contains(l, "ACK") {
				syn += 1
			}
			if strings.Contains(l, "LISTEN") {
				listen += 1
			}
		}
		m["timewait"] = timewait
		m["estab"] = estab
		m["close"] = closed
		m["syn"] = syn
		m["listen"] = listen
	}
	return m, nil
}
func NetstatMetrics() (L []*model.MetricValue) {
	tcpstatus, err := GetNetStat()
	if err != nil {
		log.Info("Get netstat error", err)
		return
	}

	for key, value := range tcpstatus {
		L = append(L, GaugeValue("netstat."+key, value))
	}
	return
}
