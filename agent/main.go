package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/anchnet/service-monitor/agent/cron"
	"github.com/anchnet/service-monitor/agent/funcs"
	"github.com/anchnet/service-monitor/agent/g"
	"github.com/anchnet/service-monitor/agent/http"
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

	if *check {
		funcs.CheckCollector()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	//init seelog
	g.InitSeeLog()

	g.InitRootDir()
	g.InitLocalIps()
	g.InitRpcClients()

	funcs.BuildMappers()

	go cron.InitDataHistory()

	cron.ReportAgentStatus()
	cron.SyncMinePlugins()
	cron.SyncBuiltinMetrics()
	cron.SyncTrustableIps()

	ReportSysInfo()

	cron.Collect()

	go http.Start()

	select {}

}
