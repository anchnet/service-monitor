package funcs

import (
	"fmt"
	"github.com/open-falcon/common/model"
	"log"
	"os/exec"
	"strings"
	"time"
)

func UptimeMetrics() []*model.MetricValue {

	out, _ := exec.Command("wmic", "os", "get", "LastBootUpTime").Output()
	boot := strings.Split(string(out[:]), "\r\n")[1]
	boot = strings.TrimSpace(boot)

	loc, err := time.LoadLocation("Local")
	if err != nil {
		log.Println(err)
	}

	year := boot[0:4]
	month := boot[4:6]
	day := boot[6:8]
	h := boot[8:10]
	m := boot[10:12]
	s := boot[12:14]

	boottime := fmt.Sprintf("%s-%s-%sT%s:%s:%sZ", year, month, day, h, m, s)
	bootDate, err := time.ParseInLocation(time.RFC3339, boottime, loc)
	if err != nil {
		log.Println(err)
	}

	out, _ = exec.Command("wmic", "os", "get", "LocalDateTime").Output()
	now := strings.Split(string(out), "\r\n")[1]

	year = now[0:4]
	month = now[4:6]
	day = now[6:8]
	h = now[8:10]
	m = now[10:12]
	s = now[12:14]
	nowtime := fmt.Sprintf("%s-%s-%sT%s:%s:%sZ", year, month, day, h, m, s)
	nowDate, err := time.ParseInLocation(time.RFC3339, nowtime, loc)
	if err != nil {
		log.Println(err)
	}

	dlt := nowDate.Sub(bootDate)
	v := dlt.Seconds()

	return []*model.MetricValue{GaugeValue("uptime", v)}
}
