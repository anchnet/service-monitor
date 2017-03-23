package funcs

import (
	"log"

	"github.com/51idc/service-monitor/mongo-monitor/g"
	"github.com/open-falcon/common/model"
)

func MongoMetrics() (L []*model.MetricValue) {
	if !g.Config().Mongo.Enabled {
		log.Println("Mongo Monitor is disabled")
		return
	}
	Addr := g.Config().Mongo.Addr
	Username := g.Config().Mongo.Username
	Password := g.Config().Mongo.Password
	Authdb := g.Config().Mongo.Authdb

	serverStatus, err := mongo_serverStatus(Addr, Authdb, Username, Password)
	if err != nil {
		log.Println(err)
		return
	}

	debug := g.Config().Debug
	smartAPI_url := g.Config().SmartAPI.Url

	if g.Config().SmartAPI.Enabled {
		endpoint, err := g.Hostname()
		version := mongo_version(serverStatus)
		if err == nil {
			smartAPI_Push(smartAPI_url, endpoint, version, debug)
		} else {
			log.Println(err)
		}
	}
	CounterMetrics, GaugeMetrics := mongo_Metrics(serverStatus)
	if connections_current, ok := GaugeMetrics["connections_current"]; ok {
		if connections_available, ok := GaugeMetrics["connections_available"]; ok {
			if connections_available != 0 {
				connections_used_percent := 100 * float64(connections_current) / float64(connections_available)
				L = append(L, GaugeValue("Mongo.connections_used_percent", connections_used_percent))
			}
		}
	}
	for metric, value := range GaugeMetrics {
		L = append(L, GaugeValue("Mongo." + metric, value))
	}
	for metric, value := range CounterMetrics {
		L = append(L, CounterValue("Mongo." + metric, value))
	}
	return
}
