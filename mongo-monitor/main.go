package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/anchnet/service-monitor/mongo-monitor/cron"
	"github.com/anchnet/service-monitor/mongo-monitor/funcs"
	"github.com/anchnet/service-monitor/mongo-monitor/g"
	"github.com/anchnet/service-monitor/mongo-monitor/http"
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

	//g.InitRootDir()
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
