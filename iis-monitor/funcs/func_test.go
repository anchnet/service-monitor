package funcs

import (
	"testing"
)

const (
	pushurl = "http://127.0.0.1/api/sevice/version"
)

func Test_iis_status(t *testing.T) {
	result, err := iis_status("_Total", "Total Bytes Received")
	t.Error(err)
	t.Log(result)
}

func Test_iis_version(t *testing.T) {
	result, err := iis_version()
	t.Error(err)
	t.Log(result)
}

func Test_smartAPI_Push(t *testing.T) {
	endpoint := "test"
	version := "1.1.1"
	smartAPI_Push(pushurl, endpoint, version, true)
}
