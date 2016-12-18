package funcs

import (
	sigar "github.com/elastic/gosigar"
	"github.com/toolkits/nux"
	"testing"
)

func Test_ProcPrecents(t *testing.T) {
	psusage_name, _, err := ProcPrecents()
	ps, _ := nux.AllProcs()
	mem := sigar.ProcMem{}
	for _, proc := range ps {
		if proc.Name == "falcon-swcollec" {
			mem.Get(proc.Pid)
			t.Log(mem.Size)
			t.Log(mem.Resident)
			t.Log(mem.Share)
		}
	}
	t.Log(psusage_name["falcon-swcollec"])
	t.Error(err)
}
