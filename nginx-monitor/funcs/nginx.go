package funcs

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/51idc/service-monitor/nginx-monitor/g"
	"github.com/open-falcon/common/model"
	"github.com/toolkits/file"
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

func nginx_version() (string, error) {
	cmd := exec.Command("nginx", "-v")
	out, err := cmd.CombinedOutput() //输出stdout和stderr的结果，此处有坑，nginx -v 的输出在stderr里，如果只打印stdout会为空
	if err != nil {
		return "", err
	}
	reader := bufio.NewReader(bytes.NewBuffer(out))
	line, err := file.ReadLine(reader)
	if err != nil {
		return "", err
	}
	v := strings.Split(string(line), ": ")
	version := v[1]
	return version, err

}

func pid_uptime(pid string) (int64, error) {
	cmd := exec.Command("stat", "-c", "%Y", pid)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	reader := bufio.NewReader(bytes.NewBuffer(out))
	line, err := file.ReadLine(reader)
	if err != nil {
		return 0, err
	}
	uptime, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, err
	}
	timenow := time.Now().Unix()
	uptime = timenow - uptime
	return uptime, err
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
	defer func() {
		if r := recover(); r != nil {
			log.Println("Nginx Recovered in Panic", r)
		}
	}()

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
	pid := g.Config().Nginx.Pid
	debug := g.Config().Debug
	smartAPI_url := g.Config().SmartAPI.Url

	if g.Config().SmartAPI.Enabled {
		endpoint, err := g.Hostname()
		version, err := nginx_version()
		if err == nil {
			smartAPI_Push(smartAPI_url, endpoint, version, debug)
		} else {
			log.Println(err)
		}
	}

	uptime, err := pid_uptime(pid)
	if err != nil {
		log.Println(err)
	} else {
		L = append(L, GaugeValue("Nginx.Uptime", uptime))
	}

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
