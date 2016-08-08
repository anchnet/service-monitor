package funcs

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/toolkits/file"

	"github.com/51idc/service-monitor/iis-monitor/g"
	"github.com/open-falcon/common/model"
)

const (
	GUAGE   = 0
	COUNTER = 1
)

type iis struct {
	metric string
	value  float64
	Type   string
	Tag    string
}

func iis_status(site string, counter string) (float64, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return 0, err
	}
	Counter := `\\` + hostname + `\web service(` + site + `)\` + counter
	value, err := ReadPerformanceCounter(Counter)
	return value, err
}

func iis_status_metric(site string, counter string, Type string, ch chan iis) {
	var IIS iis
	value, err := iis_status(site, counter)
	if err != nil {
		log.Println(err)
	}
	metric := strings.Replace(counter, " ", "_", -1)
	site_tag := strings.Replace(site, " ", "_", -1)
	IIS.metric = metric
	IIS.Type = Type
	IIS.Tag = site_tag
	IIS.value = value
	if counter == "Service Uptime" {
		IIS.metric = "Uptime"
	}
	ch <- IIS
}

func iis_version() (string, error) {
	cmd := exec.Command("powershell", "scrips/get_iis_version.ps1")
	out, err := cmd.Output()
	if err != nil {
		reader := bufio.NewReader(bytes.NewBuffer(out))
		line, _ := file.ReadLine(reader)
		return string(line), err
	}

	reader := bufio.NewReader(bytes.NewBuffer(out))
	line, err := file.ReadLine(reader)
	if err != nil {
		return "", err
	}
	return string(line), err
}

func iisMetrics() (L []*model.MetricValue) {
	if !g.Config().IIs.Enabled {
		g.Logger().Println("IIs Monitor is disabled")
		return
	}
	websites := g.Config().IIs.Websites
	debug := g.Config().Debug
	smartAPI_url := g.Config().SmartAPI.Url

	if g.Config().SmartAPI.Enabled {
		result, err := iis_version()
		endpoint, _ := g.Hostname()
		if err == nil {
			version := result
			smartAPI_Push(smartAPI_url, endpoint, version, debug)
		} else {
			g.Logger().Println(err, result)
		}
	}

	websites = append(websites, "_Total")
	chs := make([]chan iis, len(websites)*len(Mertics))
	startTime := time.Now()
	i := 0
	for _, site := range websites {
		for metric, Type := range Mertics {
			chs[i] = make(chan iis)
			switch Type {
			case "GUAGE":
				go iis_status_metric(site, metric, Type, chs[i])
			case "COUNTER":
				go iis_status_metric(site, metric, Type, chs[i])
			}
			i = i + 1
		}
	}
	for _, ch := range chs {
		IIS := <-ch
		if IIS.Type == "GUAGE" {
			L = append(L, GaugeValue("IIs."+IIS.metric, IIS.value, "site="+IIS.Tag))
		}
		if IIS.Type == "COUNTER" {
			L = append(L, CounterValue("IIs."+IIS.metric, IIS.value, "site="+IIS.Tag))
		}
	}
	endTime := time.Now()
	g.Logger().Println("IIsStats Collect complete. Process time %s.", endTime.Sub(startTime))
	return
}
