package http

import (
	"net/http"

	"github.com/51idc/service-monitor/agent/funcs"
)

func configIoStatRoutes() {
	http.HandleFunc("/page/diskio", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.IOStatsForPage())
	})
}
