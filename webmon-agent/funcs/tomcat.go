package funcs

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/51idc/service-monitor/webmon-agent/g"
	"github.com/open-falcon/common/model"
)

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

	respbody, resp_code, err := TomcathttpGet(username, password, url)
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
			L = append(L, GaugeValue("Tomcat.Connector.ThreadInfo.MaxThreads", v.ThreadInfo.MaxThreads, "Connector="+v.Name))
			L = append(L, GaugeValue("Tomcat.Connector.ThreadInfo.CurrentThreadCount", v.ThreadInfo.CurrentThreadCount, "Connector="+v.Name))
			L = append(L, GaugeValue("Tomcat.Connector.ThreadInfo.CurrentThreadsBusy", v.ThreadInfo.CurrentThreadsBusy, "Connector="+v.Name))
			L = append(L, GaugeValue("Tomcat.Connector.RequestInfo.MaxTime", v.RequestInfo.MaxTime, "Connector="+v.Name))
			L = append(L, GaugeValue("Tomcat.Connector.RequestInfo.ProcessingTime", v.RequestInfo.ProcessingTime, "Connector="+v.Name))
			L = append(L, CounterValue("Tomcat.Connector.RequestInfo.RequestCount", v.RequestInfo.RequestCount, "Connector="+v.Name))
			L = append(L, CounterValue("Tomcat.Connector.RequestInfo.ErrorCount", v.RequestInfo.ErrorCount, "Connector="+v.Name))
			L = append(L, CounterValue("Tomcat.Connector.RequestInfo.BytesReceived", v.RequestInfo.BytesReceived, "Connector="+v.Name))
			L = append(L, CounterValue("Tomcat.Connector.RequestInfo.BytesSent", v.RequestInfo.BytesSent, "Connector="+v.Name))
		}
	}
	return
}
