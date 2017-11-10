package funcs

import (
	log "github.com/cihub/seelog"
	"strconv"
	"strings"

	"github.com/anchnet/service-monitor/agent/g"
	"github.com/open-falcon/common/model"
	"github.com/toolkits/sys"
)

func DuMetrics() (L []*model.MetricValue) {
	paths := g.DuPaths()
	for _, path := range paths {
		out, err := sys.CmdOutNoLn("du", "-bs", path)
		if err != nil {
			log.Info("du -bs", path, "fail", err)
			continue
		}

		arr := strings.Fields(out)
		if len(arr) == 1 {
			continue
		}

		size, err := strconv.ParseUint(arr[0], 10, 64)
		if err != nil {
			log.Info("cannot parse du -bs", path, "output")
			continue
		}

		L = append(L, GaugeValue(g.DU_BS, size, "path="+path))
	}

	return
}
