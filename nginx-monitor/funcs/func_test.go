package funcs

import (
	"testing"
)

const (
	pushurl   = "http://127.0.0.1/api/sevice/version"
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

func Test_SendData(t *testing.T) {
	var data smartAPI_Data
	data.Endpoint = "qfeng-pc"
	data.Version = "1.1.1"
	res, err := sendData(pushurl, data)
	t.Log("res: ", res)
	t.Error(err)
}
