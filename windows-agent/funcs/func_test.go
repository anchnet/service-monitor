package funcs

import (
	"testing"
	//	sigar "github.com/elastic/gosigar"
	//	"github.com/shirou/gopsutil/process"
)

/*
func Test_ProcPercent(t *testing.T) {
	a, b, err := ProcPrecents()
	if err != nil {
		t.Error(err)
	}
	t.Log(a)
	t.Log(len(a))
	t.Log(b)
	t.Log(len(b))
	t.Log(a["QQ.exe"])
}

func Test_Processes(t *testing.T) {
	p, err := process.NewProcess(2504)
	m, err := p.MemoryInfo()
	t.Log(err)
	t.Log(m)
}
*/
func Test_NetStat(t *testing.T) {
	netstat, err := GetNetStat()

	t.Log(netstat)
	t.Error(err)
}
