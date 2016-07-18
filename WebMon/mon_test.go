package main

import (
	"testing"
)

const (
	nginx_url  = "http://127.0.0.1/status"
	apache_url = "https://www.apache.org/server-status?auto"
	tomcat_url = "http://127.0.0.1:8080/manager/status?XML=true"
	username   = "admin"
	password   = "manager"
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

func Test_tomcat(t *testing.T) {
	if text, code, err := TomcathttpGet(username, password, tomcat_url); err != nil {
		t.Error(err)
	} else {
		t.Log("text :", text)
		t.Log("code:", code)
		tomcat, err := xml_struct(text)
		t.Log("tomcat:", tomcat)
		t.Error(err)
		data := tomcat_data(tomcat)
		t.Log("data:", data)
	}
}

