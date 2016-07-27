package funcs

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/51idc/service-monitor/apache-monitor/g"
)

type smartAPI_Data struct {
	Endpoint string
	Version  string
}

func sendData(url string, data smartAPI_Data) (int, error) {

	js, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	res, err := http.Post(url, "Content-Type: application/json", bytes.NewBuffer(js))
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	return res.StatusCode, err
}

func smartAPI_Push(url string, version string) {
	var data smartAPI_Data
	endpoint, err := g.Hostname()
	if err != nil {
		log.Println(err)
		return
	}
	data.Endpoint = endpoint
	data.Version = version
	res, err := sendData(g.Config().SmartAPI.Url, data)
	if err != nil {
		log.Println(err)
		return
	}
	if res != 200 {
		log.Println("smartAPI error,statcode= ", res)
		return
	}
	if g.Config().Debug {
		log.Println("Version: ", version)
	}
	return
}
