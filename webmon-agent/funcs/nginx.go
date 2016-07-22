package funcs

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/51idc/service-monitor/webmon-agent/g"
	"github.com/open-falcon/common/model"
)

type NginxStatus struct {
	ActiveConn     int
	ServerAccepts  uint64
	ServerHandled  uint64
	ServerRequests uint64
	ServerReading  int
	ServerWriting  int
	ServerWaiting  int
}

func httpGet(url string) (string, int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	return string(body), resp.StatusCode, err
}

func nginx_status(body string) (NginxStatus, error) {
	var status NginxStatus
	str := strings.Split(body, "\n")
	s := strings.Split(str[0], " ")[2]
	ActiveConn, err := strconv.Atoi(s)
	s = strings.Split(str[2], " ")[1]
	ServerAccepts, err := strconv.ParseUint(s, 10, 64)
	s = strings.Split(str[2], " ")[2]
	ServerHandled, err := strconv.ParseUint(s, 10, 64)
	s = strings.Split(str[2], " ")[3]
	ServerRequests, err := strconv.ParseUint(s, 10, 64)
	s = strings.Split(str[3], " ")[1]
	ServerReading, err := strconv.Atoi(s)
	s = strings.Split(str[3], " ")[3]
	ServerWriting, err := strconv.Atoi(s)
	s = strings.Split(str[3], " ")[5]
	ServerWaiting, err := strconv.Atoi(s)
	status.ActiveConn = ActiveConn
	status.ServerAccepts = ServerAccepts
	status.ServerHandled = ServerHandled
	status.ServerReading = ServerReading
	status.ServerRequests = ServerRequests
	status.ServerWaiting = ServerWaiting
	status.ServerWriting = ServerWriting
	return status, err
}

func NginxMetrics() (L []*model.MetricValue) {
	if !g.Config().Nginx.Enabled {
		log.Println("Nginx Monitor is disbaled")
		return
	}
	url := g.Config().Nginx.Staturl
	respbody, resp_code, err := httpGet(url)
	if err != nil {
		log.Println(err)
		return
	}
	if resp_code != 200 {
		log.Println("Http Statu Page Open Error")
		return
	}
	stat, err := nginx_status(respbody)
	if err != nil {
		log.Println(err)
		return
	}
	L = append(L, GaugeValue("Nginx.ActiveConn", stat.ActiveConn))
	L = append(L, CounterValue("Nginx.ServerAccepts", stat.ServerAccepts))
	L = append(L, CounterValue("Nginx.ServerHandled", stat.ServerHandled))
	L = append(L, CounterValue("Nginx.ServerRequests", stat.ServerRequests))
	L = append(L, GaugeValue("Nginx.ServerReading", stat.ServerReading))
	L = append(L, GaugeValue("Nginx.ServerWaiting", stat.ServerWaiting))
	L = append(L, GaugeValue("Nginx.ServerWriting", stat.ServerWriting))
	return
}
