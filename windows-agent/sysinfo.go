package main

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/cihub/seelog"
	"net/http"

	"github.com/anchnet/service-monitor/windows-agent/g"

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
				log.Info("report sysinfo success")
				break
			}
			log.Info("report sysinfo fail")
			time.Sleep(time.Minute)
		}
	}()
}

func reportSysInfo() error {
	sysInfoInit()//将cmd代码页编码格式修改为utf8，防止中文乱码
	cpu := getCpuInfo()
	mem := getMemInfo()
	kernel := getKernelInfo()
	osVersion := getOSVersion()
	hostname, err := g.Hostname()
	if err != nil {
		log.Info(err)
		hostname = ""
	}
	sysinfo := SysInfo{
		Endpoint:  hostname,
		Cpu:       cpu,
		Mem:       mem,
		Version:   kernel,
		OSVersion: osVersion,
	}
	if g.Config().Debug {
		log.Info("sysinfo report: ", sysinfo)
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

func sysInfoInit() {
	err := exec.Command("chcp", "65001").Run()
	if err != nil {
		log.Info(err)
	}
}

func getCpuInfo() int {

	out, err := exec.Command("wmic", "cpu", "get", "numberofcores").Output()
	if err != nil {
		log.Info(err)
		return 0
	}
	line := strings.Split(string(out), "\r\n")[1]
	cores, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		log.Info(err)
	}
	return cores
}

func getMemInfo() int {
	out, err := exec.Command("wmic", "memorychip", "get", "Capacity").Output()
	if err != nil {
		log.Info(err)
		return 0
	}
	line := strings.Split(string(out), "\r\n")[1]
	mem, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		log.Info(err)
		return 0
	}
	return mem / 1024.0
}

func getKernelInfo() string {
	out, err := exec.Command("systeminfo").Output()
	if err != nil {
		log.Info(err)
		return ""
	}
	kernel := strings.Split(string(out[:]), "\r\n")[3]
	r := strings.TrimSpace(strings.Split(kernel, ":")[1])
	return r
}

func getOSVersion() string {
	out, err := exec.Command("systeminfo").Output()
	if err != nil {
		log.Info(err)
		return ""
	}
	osVersion := strings.Split(string(out[:]), "\r\n")[2]
	r := strings.TrimSpace(strings.Split(osVersion, ":")[1])
	return r
}
