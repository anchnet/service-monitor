package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/anchnet/service-monitor/windows-agent/cron"
	"github.com/anchnet/service-monitor/windows-agent/funcs"
	"github.com/anchnet/service-monitor/windows-agent/g"
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

	//init seelog
	g.InitSeeLog()

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
