package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/cihub/seelog"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/anchnet/service-monitor/agent/g"
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

func getCpuInfo() int {

	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		log.Info(err)
		return 0
	}
	defer file.Close()

	cores := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, "processor") {
			cores = cores + 1
		}
	}
	if err := scanner.Err(); err != nil {
		log.Info(err)
		return 0
	}
	return cores
}

func getMemInfo() int {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Info(err)
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, "MemTotal") {
			memTotal, _ := regexp.Compile(`\d+`)
			m := memTotal.FindString(s)
			i, _ := strconv.Atoi(m)
			return i
		}
	}
	if err := scanner.Err(); err != nil {
		log.Info(err)
		return 0
	}
	return 0
}

func getKernelInfo() string {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		log.Info(err)
		return ""
	}
	return string(out[:])
}

func getOSVersion() string {
	out, err := exec.Command("cat", "/etc/issue").Output()
	if err != nil {
		log.Info(err)
		return ""
	}
	return string(out[:])
}
