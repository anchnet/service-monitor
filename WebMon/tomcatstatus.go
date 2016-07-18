package main

import (
	"encoding/xml"
	"io/ioutil"

	"net/http"
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

func tomcat_data(status Tomcat) []*MetaData {
	data := make([]*MetaData, 0)
	Tomcat_Value := NewMetric("Tomcat.Jvm.Memory.Free")
	Tomcat_Value.SetValue(status.Jvm.Memory.Free)
	data = append(data, Tomcat_Value)
	Tomcat_Value = NewMetric("Tomcat.Jvm.Memory.Total")
	Tomcat_Value.SetValue(status.Jvm.Memory.Total)
	data = append(data, Tomcat_Value)
	Tomcat_Value = NewMetric("Tomcat.Jvm.Memory.Max")
	Tomcat_Value.SetValue(status.Jvm.Memory.Max)
	data = append(data, Tomcat_Value)
	Tomcat_Value = NewMetric("Tomcat.Jvm.Memory.usage")
	Tomcat_Jvm_Memory_usage := float64(status.Jvm.Memory.Total-status.Jvm.Memory.Free) / float64(status.Jvm.Memory.Total)
	Tomcat_Value.SetValue(int(Tomcat_Jvm_Memory_usage * 100))
	data = append(data, Tomcat_Value)
	if status.Jvm.Memorypool != nil {
		for _, v := range status.Jvm.Memorypool {
			Tomcat_Value := NewMetric("Tomcat.Jvm.Memorypool.Initial")
			Tomcat_Value.SetValue(v.UsageInit)
			Tomcat_Value.SetTags("Name=" + v.Name + ",Type=" + v.Type)
			data = append(data, Tomcat_Value)
			Tomcat_Value = NewMetric("Tomcat.Memorypool.Committed")
			Tomcat_Value.SetValue(v.UsageCommitted)
			Tomcat_Value.SetTags("Name=" + v.Name + ",Type=" + v.Type)
			data = append(data, Tomcat_Value)
			Tomcat_Value = NewMetric("Tomcat.Memorypool.Max")
			Tomcat_Value.SetValue(v.UsageMax)
			Tomcat_Value.SetTags("Name=" + v.Name + ",Type=" + v.Type)
			data = append(data, Tomcat_Value)
			Tomcat_Value = NewMetric("Tomcat.Memorypool.Used")
			Tomcat_Value.SetValue(v.UsageUsed)
			Tomcat_Value.SetTags("Name=" + v.Name + ",Type=" + v.Type)
			data = append(data, Tomcat_Value)
			Tomcat_Value = NewMetric("Tomcat.Memorypool.Usage")
			if v.UsageMax > 0 {
				Tomcat_Jvm_Memorypool_usage := float64(v.UsageUsed) / float64(v.UsageMax)
				Tomcat_Value.SetValue(int(Tomcat_Jvm_Memorypool_usage * 100))
			} else {
				Tomcat_Value.SetValue(0)
			}
			Tomcat_Value.SetTags("Name=" + v.Name + ",Type=" + v.Type)
			data = append(data, Tomcat_Value)
		}
	}
	if status.Connector != nil {
		for _, v := range status.Connector {
			Tomcat_Value := NewMetric("Tomcat.Connector.ThreadInfo.MaxThreads")
			Tomcat_Value.SetValue(v.ThreadInfo.MaxThreads)
			Tomcat_Value.SetTags("Connector=" + v.Name)
			data = append(data, Tomcat_Value)
			Tomcat_Value = NewMetric("Tomcat.Connector.ThreadInfo.CurrentThreadCount")
			Tomcat_Value.SetValue(v.ThreadInfo.CurrentThreadCount)
			Tomcat_Value.SetTags("Connector=" + v.Name)
			data = append(data, Tomcat_Value)
			Tomcat_Value = NewMetric("Tomcat.Connector.ThreadInfo.CurrentThreadsBusy")
			Tomcat_Value.SetValue(v.ThreadInfo.CurrentThreadsBusy)
			Tomcat_Value.SetTags("Connector=" + v.Name)
			data = append(data, Tomcat_Value)
			Tomcat_Value = NewMetric("Tomcat.Connector.RequestInfo.MaxTime")
			Tomcat_Value.SetValue(v.RequestInfo.MaxTime)
			Tomcat_Value.SetTags("Connector=" + v.Name)
			data = append(data, Tomcat_Value)
			Tomcat_Value = NewMetric("Tomcat.Connector.RequestInfo.ProcessingTime")
			Tomcat_Value.SetValue(v.RequestInfo.ProcessingTime)
			Tomcat_Value.SetTags("Connector=" + v.Name)
			data = append(data, Tomcat_Value)
			Tomcat_Value = NewMetric("Tomcat.Connector.RequestInfo.ProcessingTime")
			Tomcat_Value.SetValue(v.RequestInfo.ProcessingTime)
			Tomcat_Value.SetTags("Connector=" + v.Name)
			data = append(data, Tomcat_Value)
			Tomcat_Value = NewMetric("Tomcat.Connector.RequestInfo.RequestCount")
			Tomcat_Value.SetValue(v.RequestInfo.RequestCount)
			Tomcat_Value.SetTags("Connector=" + v.Name)
			data = append(data, Tomcat_Value)
			Tomcat_Value = NewMetric("Tomcat.Connector.RequestInfo.ErrorCount")
			Tomcat_Value.SetValue(v.RequestInfo.ErrorCount)
			Tomcat_Value.SetTags("Connector=" + v.Name)
			data = append(data, Tomcat_Value)
			Tomcat_Value = NewMetric("Tomcat.Connector.RequestInfo.BytesReceived")
			Tomcat_Value.SetValue(v.RequestInfo.BytesReceived)
			Tomcat_Value.SetTags("Connector=" + v.Name)
			data = append(data, Tomcat_Value)
			Tomcat_Value = NewMetric("Tomcat.Connector.RequestInfo.BytesSent")
			Tomcat_Value.SetValue(v.RequestInfo.BytesSent)
			Tomcat_Value.SetTags("Connector=" + v.Name)
			data = append(data, Tomcat_Value)
		}
	}
	return data
}
