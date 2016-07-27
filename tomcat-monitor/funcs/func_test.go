package funcs

import (
	"strings"
	"testing"
)

const (
	pushurl    = "http://127.0.0.1/api/sevice/version"
	tomcat_url = "http://127.0.0.1:8080/manager/status/"
	username   = "admin"
	password   = "manager"
)

func Test_tomcat(t *testing.T) {
	url := strings.Split(tomcat_url, "?")[0]
	staturl := url + "?XML=true"
	statallurl := url + "/all"
	if text, code, err := TomcathttpGet(username, password, staturl); err != nil {
		t.Error(err)
	} else {
		//t.Log("text :", text)
		t.Log("code:", code)
		tomcat, err := xml_struct(text)
		t.Log("tomcat:", tomcat)
		t.Error(err)
		version, err := tomcat_version(username, password, url)
		t.Log("version:", version)
		t.Error(err)
		uptime, err := tomcat_uptime(username, password, statallurl)
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
