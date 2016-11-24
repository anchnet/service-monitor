package funcs

import (
	"fmt"
	"log"

	"time"

	//	"github.com/51idc/service-monitor/windows-agent/g"
	sigar "github.com/elastic/gosigar"
	"github.com/open-falcon/common/model"
	"github.com/shirou/gopsutil/process"
)

type ProcUsage struct {
	CpuP   float64
	MemVms uint64
	MemRss uint64
}
type ProcTime struct {
	CpuUsage uint64
	CpuTotal uint64
}

type ps struct {
	Pid     int
	Name    string
	Cmdline string
}

func Processes() ([]ps, []int32, error) {
	var processes = []ps{}
	var PROCESS ps
	pids, err := process.Pids()
	if err != nil {
		return processes, pids, err
	}
	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err == nil {
			pname, err := p.Name()
			pcmdline, err := p.Cmdline()
			if err == nil {
				PROCESS.Name = pname
				PROCESS.Cmdline = pcmdline
				PROCESS.Pid = int(pid)
				processes = append(processes, PROCESS)
			}
		}
	}
	return processes, pids, err
}
func GetProcCpuTime(pids []int32) map[int]ProcTime {
	var CpuProcTimeMap = map[int]ProcTime{}
	var CpuProcTime ProcTime
	procCpu := sigar.ProcTime{}
	cpu := sigar.Cpu{}

	for _, pid := range pids {

		err := procCpu.Get(int(pid))
		if err != nil {
			//		if g.Config().Debug {
			//		log.Printf("error getting process cpu time for pid=%d: %v", pid, err)
			//		}
			continue
		}
		err = cpu.Get()
		if err != nil {
			log.Println("error getting cpu time for Total: ", err)
			continue
		}

		CpuProcTime.CpuUsage = procCpu.Total
		CpuProcTime.CpuTotal = cpu.Total()
		CpuProcTimeMap[int(pid)] = CpuProcTime
	}
	return CpuProcTimeMap
}
func GetProcCpuP(interval int, pids []int32) (map[int]float64, error) {

	var ProcCpuP = map[int]float64{}

	CpuProcTimeMap_1 := GetProcCpuTime(pids)
	time.Sleep(time.Duration(interval) * time.Second)
	CpuProcTimeMap_2 := GetProcCpuTime(pids)
	for pid, proccputime_2 := range CpuProcTimeMap_2 {
		if proccputime_1, ok := CpuProcTimeMap_1[pid]; ok {
			deltaTotal := proccputime_2.CpuTotal - proccputime_1.CpuTotal
			deltaUsage := proccputime_2.CpuUsage - proccputime_1.CpuUsage
			if deltaTotal == 0 {
				ProcCpuP[pid] = 0.0
			} else {
				ProcCpuP[pid] = 100 * float64(deltaUsage) / float64(deltaTotal)
			}
		}
	}

	return ProcCpuP, nil
}

func ProcPrecents() (map[string]ProcUsage, map[string]ProcUsage, error) {

	ps, pids, err := Processes()
	if err != nil {
		log.Println(err)
	}

	procusage_cmdline := map[string]ProcUsage{}
	procusage_name := map[string]ProcUsage{}
	procCpuP, err := GetProcCpuP(3, pids)
	if err != nil {
		return nil, nil, err
	}
	mem := sigar.ProcMem{}
	var cpuP float64
	for _, proc := range ps {
		pid := proc.Pid
		cmdline := proc.Cmdline
		name := proc.Name
		if cpuPrecent, ok := procCpuP[pid]; ok {
			cpuP = cpuPrecent
		} else {
			continue
		}
		if err := mem.Get(pid); err != nil {
			//	if g.Config().Debug {
			//	log.Printf("error getting process mem for pid=%d,name=%s: %v", pid, name, err)
			//		}
			continue
		}
		if proc_cmdline, ok := procusage_cmdline[cmdline]; ok {
			proc_cmdline.CpuP += cpuP
			proc_cmdline.MemVms += mem.Size
			proc_cmdline.MemRss += mem.Resident
			procusage_cmdline[cmdline] = proc_cmdline
		} else {
			proc_cmdline.CpuP = cpuP
			proc_cmdline.MemVms = mem.Size
			proc_cmdline.MemRss = mem.Resident
			procusage_cmdline[cmdline] = proc_cmdline
		}
		if proc_name, ok := procusage_name[name]; ok {
			proc_name.CpuP += cpuP
			proc_name.MemVms += mem.Size
			proc_name.MemRss += mem.Resident
			procusage_name[name] = proc_name
		} else {
			proc_name.CpuP = cpuP
			proc_name.MemVms = mem.Size
			proc_name.MemRss = mem.Resident
			procusage_name[name] = proc_name
		}
	}
	return procusage_name, procusage_cmdline, nil
}
func ProcPreMetrics() (L []*model.MetricValue) {
	startTime := time.Now()
	psusage_name, psusage_cmdline, err := ProcPrecents()
	if err != nil {
		log.Println(err)
		return
	}
	for proc_name, value := range psusage_name {
		tags := fmt.Sprintf("name=%v", proc_name)
		L = append(L, GaugeValue("proc.mem.vms", value.MemVms, tags))
		L = append(L, GaugeValue("proc.mem.rss", value.MemRss, tags))
		L = append(L, GaugeValue("proc.cpu.percentage", value.CpuP, tags))
	}
	for proc_cmdline, value := range psusage_cmdline {
		tags := fmt.Sprintf("cmdline=%v", proc_cmdline)
		L = append(L, GaugeValue("proc.mem.vms", value.MemVms, tags))
		L = append(L, GaugeValue("proc.mem.rss", value.MemRss, tags))
		L = append(L, GaugeValue("proc.cpu.percentage", value.CpuP, tags))
	}
	endTime := time.Now()
	log.Printf("UpdateProcessStats complete. Process time %s.", endTime.Sub(startTime))
	return
}
