package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/51idc/service-monitor/windows-agent/g"

	"os/exec"
	"strconv"
	"strings"
	"time"
)

type SysInfo struct {
	Endpoint  string `json:"endpoint"`
	Cpu       int    `json:"cpu"`
	Mem       int    `json:"mem"`
	Version   string `json:"version"`
	OSVersion string `json:"os_version"`
}

func ReportSysInfo() {
	go func() {
		for {
			if err := reportSysInfo(); err == nil {
				log.Println("report sysinfo success")
				break
			}
			log.Println("report sysinfo fail")
			time.Sleep(time.Minute)
		}
	}()
}

func reportSysInfo() error {
	cpu := getCpuInfo()
	mem := getMemInfo()
	kernel := getKernelInfo()
	osVersion := getOSVersion()
	hostname, _ := g.Hostname()

	sysinfo := SysInfo{
		Endpoint:  hostname,
		Cpu:       cpu,
		Mem:       mem,
		Version:   kernel,
		OSVersion: osVersion,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(sysinfo)

	res, err := http.Post(g.Config().SmartAPI, "application/json", b)
	if err != nil {
		return err
	}
	var message struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	json.NewDecoder(res.Body).Decode(&message)

	if message.Status == "error" {
		err = errors.New(message.Message)
	}
	return err
}

func getCpuInfo() int {

	out, _ := exec.Command("wmic", "cpu", "get", "numberofcores").Output()
	line := strings.Split(string(out), "\r\n")[1]
	cores, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		log.Println(err)
	}
	return cores
}

func getMemInfo() int {
	out, _ := exec.Command("wmic", "memorychip", "get", "Capacity").Output()
	line := strings.Split(string(out), "\r\n")[1]
	mem, _ := strconv.Atoi(strings.TrimSpace(line))
	return mem / 1024.0
}

func getKernelInfo() string {
	out, _ := exec.Command("systeminfo").Output()
	kernel := strings.Split(string(out[:]), "\r\n")[3]
	r := strings.TrimSpace(strings.Split(kernel, ":")[1])
	return r
}

func getOSVersion() string {
	out, _ := exec.Command("systeminfo").Output()
	osVersion := strings.Split(string(out[:]), "\r\n")[2]
	r := strings.TrimSpace(strings.Split(osVersion, ":")[1])
	return r
}
