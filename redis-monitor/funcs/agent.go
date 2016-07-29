package funcs

import (
	"github.com/open-falcon/common/model"
)

func AgentMetrics() []*model.MetricValue {
	return []*model.MetricValue{GaugeValue("Redis.Monitor.alive", 1)}
}
