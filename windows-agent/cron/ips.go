package cron

import (
	log "github.com/cihub/seelog"
	"time"

	"github.com/anchnet/service-monitor/windows-agent/g"
	"github.com/open-falcon/common/model"
)

func SyncTrustableIps() {
	if g.Config().Heartbeat.Enabled && g.Config().Heartbeat.Addr != "" {
		go syncTrustableIps()
	}
}

func syncTrustableIps() {

	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second

	for {
	REST:
		time.Sleep(duration)

		var ips string
		err := g.HbsClient.Call("Agent.TrustableIps", model.NullRpcRequest{}, &ips)
		if err != nil {
			log.Info("ERROR: call Agent.TrustableIps fail", err)
			goto REST
		}

		g.SetTrustableIps(ips)
	}
}
