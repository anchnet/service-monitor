package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/51idc/service-monitor/smartAPI-proxy/g"
	"github.com/51idc/service-monitor/smartAPI-proxy/http"
)

func main() {

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")

	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	g.InitRootDir()

	go http.Start()

	select {}

}
