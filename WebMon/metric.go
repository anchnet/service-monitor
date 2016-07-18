package main

import (
	"fmt"
	"os"
	"time"
)

const (
	TIME_OUT = 45

	GUAGE   = "GAUGE"
	COUNTER = "COUNTER"
	DELTA   = ""
)

// COUNTER: Speed per second
// GAUGE: Original, DEFAULT
var DataType = map[string]string{
	"Nginx.ActiveConn":                               GUAGE,
	"Nginx.ServerAccepts":                            COUNTER,
	"Nginx.ServerHandled":                            COUNTER,
	"Nginx.ServerRequests":                           COUNTER,
	"Nginx.ServerWaiting":                            GUAGE,
	"Nginx.ServerWriting":                            GUAGE,
	"Apache.Total_Accesses":                          GUAGE,
	"Apache.Total_kBytes":                            GUAGE,
	"Apache.CPULoad":                                 GUAGE,
	"Apache.Uptime":                                  GUAGE,
	"Apache.ReqPerSec":                               GUAGE,
	"Apache,BytesPerSec":                             GUAGE,
	"Apache.BytesPerReq":                             GUAGE,
	"Apache.BusyWorkers":                             GUAGE,
	"Apache.IdleWorkers":                             GUAGE,
	"Apache.ConnsTotal":                              GUAGE,
	"Apache.ConnsAsyncWriting":                       GUAGE,
	"Apache.ConnsAsyncKeepAlive":                     GUAGE,
	"Apache.ConnsAsyncClosing":                       GUAGE,
	"Apache.Waiting_for_Connection":                  GUAGE,
	"Apache.Starting_up":                             GUAGE,
	"Apache.Reading_Request":                         GUAGE,
	"Apache.Sending_Reply":                           GUAGE,
	"Apache.Keepalive_read":                          GUAGE,
	"Apache.DNS_Lookup":                              GUAGE,
	"Apache.Closing_connection":                      GUAGE,
	"Apache.Logging":                                 GUAGE,
	"Apache.Gracefully_Finishing":                    GUAGE,
	"Apache.Idle_Cleanup_of_worker":                  GUAGE,
	"Apache.Open_slot_with_no_current_process":       GUAGE,
	"Tomcat.Jvm.Memory.Free":                         GUAGE,
	"Tomcat.Jvm.Memory.Total":                        GUAGE,
	"Tomcat.Jvm.Memory.Max":                          GUAGE,
	"Tomcat.Jvm.Memory.usage":                        GUAGE,
	"Tomcat.Jvm.Memorypool.Initial":                  GUAGE,
	"Tomcat.Jvm.Memorypool.Committed":                GUAGE,
	"Tomcat.Jvm.Memorypool.Max":                      GUAGE,
	"Tomcat.Jvm.Memorypool.Used":                     GUAGE,
	"Tomcat.Jvm.Memorypool.Usage":                    GUAGE,
	"Tomcat.Connector.ThreadInfo.MaxThreads":         GUAGE,
	"Tomcat.Connector.ThreadInfo.CurrentThreadCount": GUAGE,
	"Tomcat.Connector.ThreadInfo.CurrentThreadsBusy": GUAGE,
	"Tomcat.Connector.RequestInfo.MaxTime":           GUAGE,
	"Tomcat.Connector.RequestInfo.ProcessingTime":    GUAGE,
	"Tomcat.Connector.RequestInfo.RequestCount":      COUNTER,
	"Tomcat.Connector.RequestInfo.ErrorCount":        COUNTER,
	"Tomcat.Connector.RequestInfo.BytesReceived":     COUNTER,
	"Tomcat.Connector.RequestInfo.BytesSent":         COUNTER,
}

func dataType(key_ string) string {
	if v, ok := DataType[key_]; ok {
		return v
	}
	return GUAGE
}

type MetaData struct {
	Metric      string      `json:"metric"`      //key
	Endpoint    string      `json:"endpoint"`    //hostname
	Value       interface{} `json:"value"`       // number or string
	CounterType string      `json:"counterType"` // GAUGE  原值   COUNTER 差值(ps)
	Tags        string      `json:"tags"`        // port=3306,k=v
	Timestamp   int64       `json:"timestamp"`
	Step        int         `json:"step"`
}

func (m *MetaData) String() string {
	s := fmt.Sprintf("MetaData Metric:%s Endpoint:%s Value:%v CounterType:%s Tags:%s Timestamp:%d Step:%d",
		m.Metric, m.Endpoint, m.Value, m.CounterType, m.Tags, m.Timestamp, m.Step)
	return s
}

func NewMetric(name string) *MetaData {
	return &MetaData{
		Metric:      name,
		Endpoint:    hostname(),
		CounterType: dataType(name),
		Tags:        "",
		Timestamp:   time.Now().Unix(),
		Step:        cfg.Interval,
	}
}

func hostname() string {
	host := cfg.Endpoint
	if host != "" {
		return host
	}
	host, err := os.Hostname()
	if err != nil {
		host = cfg.Endpoint
	}
	return host
}

func (m *MetaData) SetValue(v interface{}) {
	m.Value = v
}

func (m *MetaData) SetTags(v string) {
	m.Tags = v
}
