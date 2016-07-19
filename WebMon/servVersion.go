package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"os/exec"
	"strings"
)

func ReportVersion() {
	if cfg.NginxEnabled == 1 {
		if err := reportNginx(); err == nil {
			log.Infof("report nginx version success")
		} else {
			log.Infof("report nginx version fail", err)
		}
	}
}

func reportNginx() error {
	version, err := nginxVersion()
	if err != nil {
		return err
	}
	s := fmt.Sprintf(`{"endpoint": "%s", "version": "%s"}`, hostname(), version)
	b := bytes.NewBufferString(s)
	res, err := http.Post(cfg.SmartAPI, "application/json", b)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		err = errors.New("HTTP STATUS NOT 200")
		return err
	}

	var message struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	json.NewDecoder(res.Body).Decode(&message)

	if message.Status != "ok" {
		err = errors.New(message.Message)
	}
	return err
}

func nginxVersion() (string, error) {
	//out, err := exec.Command("ndslfkjafx", "-v").Output()
	out, err := exec.Command("nginx", "-v").CombinedOutput()
	if err != nil {
		return "", err
	}
	version := strings.TrimSpace(strings.Split(string(out), ":")[1])
	return version, nil

}
