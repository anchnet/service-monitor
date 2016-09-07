package funcs

import (
	"testing"
)

const (
	pushurl  = "http://127.0.0.1/api/sevice/version"
	Addr     = "127.0.0.1:27017"
	Username = ""
	Password = ""
	Authdb   = ""
)

func Test_mongo_stat(t *testing.T) {
	serverStatus, err := mongo_serverStatus(Addr, Authdb, Username, Password)
	t.Error(err)
	ver := mongo_version(serverStatus)
	t.Log(ver)
	CounterMetrics, GaugeMetrics := mongo_Metrics(serverStatus)
	t.Log(CounterMetrics)
	t.Log(GaugeMetrics)
}

func Test_smartAPI_Push(t *testing.T) {
	endpoint := "test"
	version := "1.1.1"
	smartAPI_Push(pushurl, endpoint, version, true)

}
