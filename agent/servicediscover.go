package main

import (
	log "github.com/cihub/seelog"
	"os/exec"
	"time"
	"net/http"
	"bytes"
	"strings"
	"encoding/json"
	"errors"
	"github.com/anchnet/service-monitor/agent/g"
)

type ServiceDiscoverInfo struct {
	Endpoint        string `json:"endpoint"`
	OS              string `json:"os"`
	ServicePortList map[string]string    `json:"service_port_list"`
}

func ServiceDiscover() {
	for {
		if err := reportServiceDiscover(); err == nil {
			log.Info("discover service success")
		} else {
			log.Error("discover service fail")
		}
		time.Sleep(time.Duration(g.Config().Discover.Cycle) * time.Second)
	}
}

func reportServiceDiscover() error {
	hostname, err := g.Hostname()
	if err != nil {
		log.Error(err)
		hostname = ""
	}
	serviceMap := getServicePortList()
	sd_info := ServiceDiscoverInfo{
		Endpoint:hostname,
		OS : "linux",
		ServicePortList:serviceMap,
	}
	log.Info(sd_info)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(sd_info)

	res, err := http.Post(g.Config().Discover.Url, "application/json", b)
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

func getServicePortList() map[string]string {
	var map_result map[string]string = make(map[string]string)

	cmd := exec.Command("bash")
	in := bytes.NewBuffer(nil)
	cmd.Stdin = in//绑定输入
	var out bytes.Buffer
	cmd.Stdout = &out //绑定输出
	go func() {
		var cmd_string string = `netstat -nultp | awk 'NR>2{printf $4 " " $7 "\n"}' | awk -F : '{printf $NF "\n"}' | awk -F '[ /]+' '{printf $1 " " $3 "\n"}' | awk '!a[$1]++{print}'`
		in.WriteString(cmd_string)//写入你的命令，可以有多行，"\n"表示回车
	}()
	err := cmd.Start()
	if err != nil {
		log.Error(err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Info("Command finished with error: %v", err)
	}

	if err != nil {
		log.Info(err)
		return map_result
	} else {
		lens := strings.Split(out.String(), "\n")
		for _, each_len := range lens {
			port_name_pair := strings.Split(each_len, " ")
			if len(port_name_pair) == 2 {
				if val, ok := map_result[port_name_pair[1]]; ok {
					map_result[port_name_pair[1]] = val + "," + port_name_pair[0]
				} else {
					map_result[port_name_pair[1]] = port_name_pair[0]
				}
			}
		}
	}
	return map_result
}