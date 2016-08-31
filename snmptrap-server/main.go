package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/51idc/service-monitor/snmptrap-server/g"
	"github.com/51idc/service-monitor/snmptrap-server/http"
	"github.com/51idc/service-monitor/snmptrap-server/redis"
	"github.com/51idc/service-monitor/snmptrap-server/trap"
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
	redis.InitRedisConnPool()
	go trap.Start()
	go http.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		redis.RedisConnPool.Close()
		os.Exit(0)
	}()

	select {}

}
