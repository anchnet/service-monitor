package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/51idc/service-monitor/windows-agent/cron"
	"github.com/51idc/service-monitor/windows-agent/funcs"
	"github.com/51idc/service-monitor/windows-agent/g"
)

func main() {

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	check := flag.Bool("check", false, "check collector")

	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	if *check {
		funcs.CheckCollector()
		os.Exit(0)
	}

	//g.InitRootDir()
	g.InitLocalIps()
	g.InitRpcClients()

	funcs.BuildMappers()

	go cron.InitDataHistory()

	cron.ReportAgentStatus()
	cron.SyncBuiltinMetrics()
	cron.SyncTrustableIps()

	ReportSysInfo()

	cron.Collect()

	select {}

}
