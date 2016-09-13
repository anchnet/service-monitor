package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type smartAPI_Data struct {
	Endpoint string `json:"endpoint"`
	Version  string `json:"version"`
}

type smartAPI_Result struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func sendData(url string, data smartAPI_Data) ([]byte, int, error) {

	js, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err
	}
	res, err := http.Post(url, "Content-Type: application/json", bytes.NewBuffer(js))
	if err != nil {
		return nil, 0, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}
	defer res.Body.Close()
	return body, res.StatusCode, err
}
