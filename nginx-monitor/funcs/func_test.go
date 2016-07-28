package funcs

import (
	"testing"
)

const (
	pushurl   = "http://127.0.0.1/api/service/version"
	nginx_url = "http://127.0.0.1/status"
	pid       = "/var/run/nginx.pid"
)

func Test_nginx(t *testing.T) {
	if text, code, err := httpGet(nginx_url); err != nil {
		t.Error(err)
	} else {
		//t.Log("text :", text)
		t.Log("code:", code)
		status, _ := nginx_status(text)
		t.Log("status:", status)
		version, err := nginx_version()
		t.Log("version:", version)
		t.Error(err)
		uptime, err := pid_uptime(pid)
		t.Log("uptime:", uptime)
		t.Error(err)
	}
}

func Test_smartAPI_Push(t *testing.T) {
	endpoint := "test"
	version := "1.1.1"
	smartAPI_Push(pushurl, endpoint, version, true)

}
