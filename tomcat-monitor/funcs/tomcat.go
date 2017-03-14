package funcs

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/51idc/service-monitor/tomcat-monitor/g"
	"github.com/open-falcon/common/model")

type Tomcat struct {
	Jvm       Jvm         `xml:"jvm"`
	Connector []Connector `xml:"connector"`
}

type Jvm struct {
	Memory     Memory       `xml:"memory"`
	Memorypool []Memorypool `xml:"memorypool"`
}

type Memorypool struct {
	Name           string `xml:"name,attr"`
	Type           string `xml:"type,attr"`
	UsageInit      uint64 `xml:"usageInit,attr"`
	UsageCommitted uint64 `xml:"usageCommitted,attr"`
	UsageMax       uint64 `xml:"usageMax,attr"`
	UsageUsed      uint64 `xml:"usageUsed,attr"`
}

type Memory struct {
	Free  uint64 `xml:"free,attr"`
	Total uint64 `xml:"total,attr"`
	Max   uint64 `xml:"max,attr"`
}
type Connector struct {
	Name        string      `xml:"name,attr"`
	ThreadInfo  ThreadInfo  `xml:"threadInfo"`
	RequestInfo RequestInfo `xml:"requestInfo"`
}
type ThreadInfo struct {
	MaxThreads         uint64 `xml:"maxThreads,attr"`
	CurrentThreadCount uint64 `xml:"currentThreadCount,attr"`
	CurrentThreadsBusy uint64 `xml:"currentThreadBusy,attr"`
}
type RequestInfo struct {
	MaxTime        uint64 `xml:"maxTime,attr"`
	ProcessingTime uint64 `xml:"processingTime,attr"`
	RequestCount   uint64 `xml:"requestCount,attr"`
	ErrorCount     uint64 `xml:"errorCount,attr"`
	BytesReceived  uint64 `xml:"bytesReceived,attr"`
	BytesSent      uint64 `xml:"bytesSent,attr"`
}

func tomcat_version(username string, password string, url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", err
	}

	d, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", err
	}
	version := d.Find("table").Eq(3).Find("td.row-center").Find("small").Eq(0).Text()
	return version, err
}

func tomcat_uptime(username string, password string, url string) (int64, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		return 0, err
	}

	d, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return 0, err
	}
	var start_time_string string
	d.Find("p").Each(func(i int, contentSelection *goquery.Selection) {
		t := contentSelection.Text()
		if strings.Contains(t, "Start time") {
			start_time_string = t
		}
	})
	start_time_string = strings.Split(start_time_string, " Startup time")[0]
	start_time_string = strings.Split(start_time_string, ": ")[1]
	start_time, err := time.Parse("Mon Jan _2 15:04:05 MST 2006", start_time_string)
	uptime := time.Now().Unix() - start_time.Unix()
	return uptime, err
}

func TomcathttpGet(username string, password string, url string) (string, int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, err
	}
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
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

func xml_struct(body string) (Tomcat, error) {
	var result Tomcat
	err := xml.Unmarshal([]byte(body), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func TomcatMetrics() (L []*model.MetricValue) {
	if !g.Config().Tomcat.Enabled {
		log.Println("Tomcat Monitor is disbaled")
		return
	}
	username := g.Config().Tomcat.Username
	password := g.Config().Tomcat.Password
	url := g.Config().Tomcat.Staturl
	url = strings.Split(url, "?")[0]
	staturl := url + "?XML=true"
	statallurl := url + "/all"
	debug := g.Config().Debug
	smartAPI_url := g.Config().SmartAPI.Url

	if g.Config().SmartAPI.Enabled {
		endpoint, err := g.Hostname()
		version, err := tomcat_version(username, password, url)
		if err == nil {
			smartAPI_Push(smartAPI_url, endpoint, version, debug)
		} else {
			log.Println(err)
		}
	}

	uptime, err := tomcat_uptime(username, password, statallurl)
	if err != nil {
		log.Println(err)
	} else {
		L = append(L, GaugeValue("Tomcat.Uptime", uptime))
	}

	respbody, resp_code, err := TomcathttpGet(username, password, staturl)
	if err != nil {
		log.Println(err)
		return
	}
	if resp_code != 200 {
		log.Println("Http Statu Page Open Error")
		return
	}
	stat, err := xml_struct(respbody)
	if err != nil {
		log.Println(err)
		return
	}
	L = append(L, GaugeValue("Tomcat.Jvm.Memory.Free", stat.Jvm.Memory.Free))
	L = append(L, GaugeValue("Tomcat.Jvm.Memory.Total", stat.Jvm.Memory.Total))
	L = append(L, GaugeValue("Tomcat.Jvm.Memory.Max", stat.Jvm.Memory.Max))
	Tomcat_Jvm_Memory_usage := float64(stat.Jvm.Memory.Total-stat.Jvm.Memory.Free) / float64(stat.Jvm.Memory.Total)
	L = append(L, GaugeValue("Tomcat.Jvm.Memory.usage", int(Tomcat_Jvm_Memory_usage*100)))
	if stat.Jvm.Memorypool != nil {
		for _, v := range stat.Jvm.Memorypool {
			L = append(L, GaugeValue("Tomcat.Jvm.Memorypool.Initial", v.UsageInit, "Name="+v.Name+",Type="+v.Type))
			L = append(L, GaugeValue("Tomcat.Memorypool.Committed", v.UsageCommitted, "Name="+v.Name+",Type="+v.Type))
			L = append(L, GaugeValue("Tomcat.Memorypool.Max", v.UsageMax, "Name="+v.Name+",Type="+v.Type))
			L = append(L, GaugeValue("Tomcat.Memorypool.Used", v.UsageUsed, "Name="+v.Name+",Type="+v.Type))
			if v.UsageMax > 0 {
				Tomcat_Jvm_Memorypool_usage := float64(v.UsageUsed) / float64(v.UsageMax)
				L = append(L, GaugeValue("Tomcat.Memorypool.Usage", int(Tomcat_Jvm_Memorypool_usage*100), "Name="+v.Name+",Type="+v.Type))
			} else {
				L = append(L, GaugeValue("Tomcat.Memorypool.Usage", 0, "Name="+v.Name+",Type="+v.Type))
			}
		}
	}
	if stat.Connector != nil {
		for _, v := range stat.Connector {
			Connector := strings.Trim(v.Name, `"`)
			L = append(L, GaugeValue("Tomcat.Connector.ThreadInfo.MaxThreads", v.ThreadInfo.MaxThreads, "Connector="+Connector))
			L = append(L, GaugeValue("Tomcat.Connector.ThreadInfo.CurrentThreadCount", v.ThreadInfo.CurrentThreadCount, "Connector="+Connector))
			L = append(L, GaugeValue("Tomcat.Connector.ThreadInfo.CurrentThreadsBusy", v.ThreadInfo.CurrentThreadsBusy, "Connector="+Connector))
			L = append(L, GaugeValue("Tomcat.Connector.RequestInfo.MaxTime", v.RequestInfo.MaxTime, "Connector="+Connector))
			L = append(L, GaugeValue("Tomcat.Connector.RequestInfo.ProcessingTime", v.RequestInfo.ProcessingTime, "Connector="+Connector))
			L = append(L, CounterValue("Tomcat.Connector.RequestInfo.RequestCount", v.RequestInfo.RequestCount, "Connector="+Connector))
			L = append(L, CounterValue("Tomcat.Connector.RequestInfo.ErrorCount", v.RequestInfo.ErrorCount, "Connector="+Connector))
			L = append(L, CounterValue("Tomcat.Connector.RequestInfo.BytesReceived", v.RequestInfo.BytesReceived, "Connector="+Connector))
			L = append(L, CounterValue("Tomcat.Connector.RequestInfo.BytesSent", v.RequestInfo.BytesSent, "Connector="+Connector))
		}
	}

	return
}
