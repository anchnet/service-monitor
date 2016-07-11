package main

import (
	"fmt"
	"os"
	"time"
)

const (
	TIME_OUT = 5

	GUAGE   = "GAUGE"
	COUNTER = "COUNTER"
	DELTA   = ""
)

// COUNTER: Speed per second
// GAUGE: Original, DEFAULT
var DataType = map[string]string{
	"Nginx.ActiveConn":     GUAGE,
	"Nginx.ServerAccepts":  COUNTER,
	"Nginx.ServerHandled":  COUNTER,
	"Nginx.ServerRequests": COUNTER,
	"Nginx.ServerWaiting":  GUAGE,
	"Nginx.ServerWriting":  GUAGE,
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
