package funcs

import (
	"log"
	"strconv"
	"strings"

	"github.com/51idc/service-monitor/webmon-agent/g"
	"github.com/open-falcon/common/model"
)

func apache_status(body string) (map[string]float64, error) {
	ApacheStatus := make(map[string]float64)
	var err error
	var value float64
	str := strings.Split(body, "\n")
	for _, line := range str {
		v := strings.Split(line, ": ")
		if v[0] == "" {
			return ApacheStatus, err
		}
		if v[0] == "Scoreboard" {
			ApacheStatus["Waiting_for_Connection"] = float64(strings.Count(v[1], "_"))
			ApacheStatus["Starting_up"] = float64(strings.Count(v[1], "S"))
			ApacheStatus["Reading_Request"] = float64(strings.Count(v[1], "R"))
			ApacheStatus["Sending_Reply"] = float64(strings.Count(v[1], "W"))
			ApacheStatus["Keepalive_read"] = float64(strings.Count(v[1], "K"))
			ApacheStatus["DNS_Lookup"] = float64(strings.Count(v[1], "D"))
			ApacheStatus["Closing_connection"] = float64(strings.Count(v[1], "C"))
			ApacheStatus["Logging"] = float64(strings.Count(v[1], "L"))
			ApacheStatus["Gracefully_Finishing"] = float64(strings.Count(v[1], "G"))
			ApacheStatus["Idle_Cleanup_of_worker"] = float64(strings.Count(v[1], "I"))
			ApacheStatus["Open_slot_with_no_current_process"] = float64(strings.Count(v[1], "."))
			continue
		}
		value, err = strconv.ParseFloat(v[1], 64)
		if err == nil {
			ApacheStatus[strings.Replace(v[0], " ", "_", -1)] = value
		}
	}

	return ApacheStatus, err
}

func ApacheMetrics() (L []*model.MetricValue) {
	if !g.Config().Apache.Enabled {
		log.Println("Apache Monitor is disabled")
		return
	}
	url := g.Config().Apache.Staturl
	respbody, resp_code, err := httpGet(url)
	if err != nil {
		log.Println(err)
		return
	}
	if resp_code != 200 {
		log.Println("Http Statu Page Open Error")
		return
	}
	stat, err := apache_status(respbody)
	if err != nil {
		log.Println(err)
		return
	}

	for index, value := range stat {
		L = append(L, GaugeValue("Apache."+index, value))
	}
	return
}
