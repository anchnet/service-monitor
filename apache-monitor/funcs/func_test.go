package funcs

import (
	"strings"
	"testing"
)

const (
	pushurl    = "http://127.0.0.1/api/sevice/version"
	apache_url = "https://www.apache.org/server-status?auto"
)

func Test_apache(t *testing.T) {
	apache_url := strings.Split(apache_url, "?")[0]
	url := apache_url + "?auto"
	if text, code, err := httpGet(url); err != nil {
		t.Error(err)
	} else {
		//t.Log("text :", text)
		t.Log("code:", code)
		status, _ := apache_status(text)
		t.Log("status:", status)
		version, err := apache_version(apache_url)
		t.Log("version:", version)
		t.Error(err)
	}
}

func Test_smartAPI_Push(t *testing.T) {
	endpoint := "test"
	version := "1.1.1"
	smartAPI_Push(pushurl, endpoint, version, true)

}
