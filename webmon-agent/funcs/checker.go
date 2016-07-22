package funcs

import (
	"fmt"
)

func CheckCollector() {

	output := make(map[string]bool)

	output["NginxMetrics"] = len(NginxMetrics()) > 0
	output["ApacheMetrics"] = len(ApacheMetrics()) > 0
	output["TomcatMetrics"] = len(TomcatMetrics()) > 0

	for k, v := range output {
		status := "fail"
		if v {
			status = "ok"
		}
		fmt.Println(k, "...", status)
	}
}
