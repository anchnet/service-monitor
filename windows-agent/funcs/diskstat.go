package funcs

import (
	"log"

	"github.com/51idc/service-monitor/windows-agent/tools/disk"
	"github.com/open-falcon/common/model"
)

func DiskIOMetrics() (L []*model.MetricValue) {

	dsList, err := disk.DiskIOCounters()
	if err != nil {
		log.Println("Get devices io fail: ", err)
		return
	}

	for _, ds := range dsList {
		device := "device=" + ds.Name

		L = append(L, CounterValue("disk.io.read_requests", ds.ReadCount, device))
		L = append(L, CounterValue("disk.io.read_bytes", ds.ReadBytes, device))
		L = append(L, CounterValue("disk.io.write_requests", ds.WriteCount, device))
		L = append(L, CounterValue("disk.io.write_bytes", ds.WriteBytes, device))
		L = append(L, CounterValue("disk.io.read_time", ds.ReadTime, device))
		L = append(L, CounterValue("disk.io.write_time", ds.WriteTime, device))
		L = append(L, CounterValue("disk.io.iotime", ds.IoTime, device))
	}
	return
}
