package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type NginxStatus struct {
	ActiveConn     int
	ServerAccepts  uint64
	ServerHandled  uint64
	ServerRequests uint64
	ServerReading  int
	ServerWriting  int
	ServerWaiting  int
}

func httpGet(url string) (string, int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	return string(body), resp.StatusCode, err
}

func nginx_status(body string) (NginxStatus, error) {
	var status NginxStatus
	str := strings.Split(body, "\n")
	s := strings.Split(str[0], " ")[2]
	ActiveConn, err := strconv.Atoi(s)
	s = strings.Split(str[2], " ")[1]
	ServerAccepts, err := strconv.ParseUint(s, 10, 64)
	s = strings.Split(str[2], " ")[2]
	ServerHandled, err := strconv.ParseUint(s, 10, 64)
	s = strings.Split(str[2], " ")[3]
	ServerRequests, err := strconv.ParseUint(s, 10, 64)
	s = strings.Split(str[3], " ")[1]
	ServerReading, err := strconv.Atoi(s)
	s = strings.Split(str[3], " ")[3]
	ServerWriting, err := strconv.Atoi(s)
	s = strings.Split(str[3], " ")[5]
	ServerWaiting, err := strconv.Atoi(s)
	status.ActiveConn = ActiveConn
	status.ServerAccepts = ServerAccepts
	status.ServerHandled = ServerHandled
	status.ServerReading = ServerReading
	status.ServerRequests = ServerRequests
	status.ServerWaiting = ServerWaiting
	status.ServerWriting = ServerWriting
	return status, err
}

func nginx_data(status NginxStatus) []*MetaData {
	data := make([]*MetaData, 0)
	stat := status
	Nginx_ACtiveConn := NewMetric("Nginx.ActiveConn")
	Nginx_ACtiveConn.SetValue(stat.ActiveConn)
	data = append(data, Nginx_ACtiveConn)
	Nginx_ServerAccepts := NewMetric("Nginx.ServerAccepts")
	Nginx_ServerAccepts.SetValue(stat.ServerAccepts)
	data = append(data, Nginx_ServerAccepts)
	Nginx_ServerHandled := NewMetric("Nginx.ServerHandled")
	Nginx_ServerHandled.SetValue(stat.ServerHandled)
	data = append(data, Nginx_ServerHandled)
	Nginx_ServerReading := NewMetric("Nginx.ServerReading")
	Nginx_ServerReading.SetValue(stat.ServerReading)
	data = append(data, Nginx_ServerReading)
	Nginx_ServerRequests := NewMetric("Nginx.ServerRequests")
	Nginx_ServerRequests.SetValue(stat.ServerRequests)
	data = append(data, Nginx_ServerRequests)
	Nginx_ServerWaiting := NewMetric("Nginx.ServerWaiting")
	Nginx_ServerWaiting.SetValue(stat.ServerWaiting)
	data = append(data, Nginx_ServerWaiting)
	Nginx_ServerWriting := NewMetric("Nginx.ServerWriting")
	Nginx_ServerWriting.SetValue(stat.ServerWriting)
	data = append(data, Nginx_ServerWriting)
	return data
}
