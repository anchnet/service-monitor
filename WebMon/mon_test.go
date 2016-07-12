package main

import (
	"testing"
)

const (
	nginx_url  = "http://127.0.0.1/status"
	apache_url = "https://www.apache.org/server-status?auto"
)

func Test_nginx(t *testing.T) {
	if text, code, err := httpGet(nginx_url); err != nil {
		t.Error(err)
	} else {
		t.Log("text :", text)
		t.Log("code:", code)
		status, _ := nginx_status(text)
		t.Log("status:", status)
		data := nginx_data(status)
		t.Log("data:", data)
	}
}

func Test_apache(t *testing.T) {
	if text, code, err := httpGet(apache_url); err != nil {
		t.Error(err)
	} else {
		t.Log("text :", text)
		t.Log("code:", code)
		status, _ := apache_status(text)
		t.Log("status:", status)
		data := apache_data(status)
		t.Log("data:", data)
	}
}
