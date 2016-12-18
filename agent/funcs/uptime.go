package funcs

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/open-falcon/common/model"
)

func UptimeMetrics() []*model.MetricValue {
	out, _ := exec.Command("cat", "/proc/uptime").Output()
	s := string(out[:])
	uptime := strings.Split(s, " ")[0]
	v, _ := strconv.ParseFloat(uptime, 64)
	return []*model.MetricValue{GaugeValue("uptime", v)}
}
