package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/anchnet/service-monitor/iis-monitor/cron"
	"github.com/anchnet/service-monitor/iis-monitor/funcs"
	"github.com/anchnet/service-monitor/iis-monitor/g"
	"github.com/anchnet/service-monitor/iis-monitor/http"
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
	// rm g.InitLog()

	g.InitRootDir()
	//g.InitRpcClients()

	if *check {
		funcs.CheckCollector()
		os.Exit(0)
	}

	funcs.BuildMappers()

	cron.Collect()

	go http.Start()

	select {}
}
