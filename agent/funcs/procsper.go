package funcs

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	sigar "github.com/elastic/gosigar"
	"github.com/open-falcon/common/model"
	"github.com/toolkits/nux"
)

type ProcUsage struct {
	CpuP     float64
	MemSize  uint64
	MemRss   uint64
	MemShare uint64
}

func GetProcCpuP() (map[int]float64, error) {
	pidCPUP := map[int]float64{}

	out, err := exec.Command("top", "-b", "-n 1").Output()
	if err != nil {
		return nil, err
	} else {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines[7:] {
			//fmt.Println(line)
			fields := strings.Fields(line)
			if len(fields) == 12 {
				pid, err := strconv.Atoi(fields[0])
				cpuP, err := strconv.ParseFloat(fields[8], 64)
				if err != nil {
					log.Println(err)
					continue
				}
				pidCPUP[pid] = cpuP
			}
		}
	}
	return pidCPUP, nil
}

func ProcPrecents() (map[string]ProcUsage, map[string]ProcUsage, error) {

	ps, err := nux.AllProcs()
	if err != nil {
		return nil, nil, err
	}
	procusage_cmdline := map[string]ProcUsage{}
	procusage_name := map[string]ProcUsage{}
	procCpuP, err := GetProcCpuP()
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
			log.Println("error getting process mem for pid=%d: %v", pid, err)
			continue
		}
		if proc_cmdline, ok := procusage_cmdline[cmdline]; ok {
			proc_cmdline.CpuP += cpuP
			proc_cmdline.MemSize += mem.Size
			proc_cmdline.MemShare += mem.Share
			proc_cmdline.MemRss += mem.Resident
			procusage_cmdline[cmdline] = proc_cmdline
		} else {
			proc_cmdline.CpuP = cpuP
			proc_cmdline.MemSize = mem.Size
			proc_cmdline.MemShare = mem.Share
			proc_cmdline.MemRss = mem.Resident
			procusage_cmdline[cmdline] = proc_cmdline
		}
		if proc_name, ok := procusage_name[name]; ok {
			proc_name.CpuP += cpuP
			proc_name.MemSize += mem.Size
			proc_name.MemShare += mem.Share
			proc_name.MemRss += mem.Resident
			procusage_name[name] = proc_name
		} else {
			proc_name.CpuP = cpuP
			proc_name.MemSize = mem.Size
			proc_name.MemShare = mem.Share
			proc_name.MemRss = mem.Resident
			procusage_name[name] = proc_name
		}
		if name == "falcon-swcollec" {
			log.Println(procusage_name[name])
		}
	}
	return procusage_name, procusage_cmdline, nil
}
func ProcPreMetrics() (L []*model.MetricValue) {
	psusage_name, psusage_cmdline, err := ProcPrecents()
	if err != nil {
		log.Println(err)
		return
	}
	for proc_name, value := range psusage_name {
		tags := fmt.Sprintf("name=%v", proc_name)
		L = append(L, GaugeValue("proc.mem.size", value.MemSize, tags))
		L = append(L, GaugeValue("proc.mem.rss", value.MemRss, tags))
		L = append(L, GaugeValue("proc.mem.share", value.MemShare, tags))
		L = append(L, GaugeValue("proc.cpu.percentage", value.CpuP, tags))
	}
	for proc_cmdline, value := range psusage_cmdline {
		tags := fmt.Sprintf("cmdline=%v", proc_cmdline)
		L = append(L, GaugeValue("proc.mem.size", value.MemSize, tags))
		L = append(L, GaugeValue("proc.mem.rss", value.MemRss, tags))
		L = append(L, GaugeValue("proc.mem.share", value.MemShare, tags))
		L = append(L, GaugeValue("proc.cpu.percentage", value.CpuP, tags))
	}
	return
}
