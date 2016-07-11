package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func sendData(data []*MetaData) (int, error) {

	js, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	log.Debugf("Send to %s, size: %d", cfg.PushUrl, len(data))
	for _, m := range data {
		log.Debugf("%s", m)
	}

	res, err := http.Post(cfg.PushUrl, "Content-Type: application/json", bytes.NewBuffer(js))
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()
	return res.StatusCode, err
}
