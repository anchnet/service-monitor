package funcs

import (
	"bufio"
	"bytes"
	"fmt"

	"os/exec"
	"strings"

	"github.com/toolkits/file"

	"github.com/51idc/service-monitor/iis-monitor/g"
	"github.com/open-falcon/common/model"
)

func in_array(a string, array []string) bool {
	for _, v := range array {
		if a == v {
			return true
		}
	}
	return false
}

func format_mertic(metric string) string {
	result := strings.TrimSpace(metric)
	result = strings.Replace(result, " ", "_", -1)
	return result
}

func iis_version() (string, error) {
	cmd := exec.Command("powershell", "scrips/get_iis_version.ps1")
	out, err := cmd.Output()
	if err != nil {
		reader := bufio.NewReader(bytes.NewBuffer(out))
		line, _ := file.ReadLine(reader)
		return string(line), err
	}

	reader := bufio.NewReader(bytes.NewBuffer(out))
	line, err := file.ReadLine(reader)
	if err != nil {
		return "", err
	}
	return string(line), err
}

func iisMetrics() (L []*model.MetricValue) {
	if !g.Config().IIs.Enabled {
		g.Logger().Println("IIs Monitor is disabled")
		return
	}
	websites := g.Config().IIs.Websites
	debug := g.Config().Debug
	smartAPI_url := g.Config().SmartAPI.Url

	if g.Config().SmartAPI.Enabled {
		result, err := iis_version()
		endpoint, _ := g.Hostname()
		if err == nil {
			version := result
			smartAPI_Push(smartAPI_url, endpoint, version, debug)
		} else {
			g.Logger().Println(err, result)
		}
	}

	websites = append(websites, "_Total")
	IIsStat, err := IIsCounters()
	if err != nil {
		g.Logger().Println(err)
		return
	}
	for _, iisStat := range IIsStat {
		if in_array(iisStat.Name, websites) {
			tag := fmt.Sprintf("site=%s", format_mertic(iisStat.Name))
			L = append(L, CounterValue("iis.bytes.received", iisStat.BytesReceivedPersec, tag))
			L = append(L, CounterValue("iis.bytes.sent", iisStat.BytesSentPersec, tag))
			L = append(L, CounterValue("iis.requests.cgi", iisStat.CGIRequestsPersec, tag))
			L = append(L, CounterValue("iis.connection.attempts", iisStat.ConnectionAttemptsPersec, tag))
			L = append(L, CounterValue("iis.requests.copy", iisStat.CopyRequestsPersec, tag))
			L = append(L, GaugeValue("iis.connections", iisStat.CurrentConnections, tag))
			L = append(L, CounterValue("iis.requests.delete", iisStat.DeleteRequestsPersec, tag))
			L = append(L, CounterValue("iis.requests.get", iisStat.GetRequestsPersec, tag))
			L = append(L, CounterValue("iis.requests.head", iisStat.HeadRequestsPersec, tag))
			L = append(L, CounterValue("iis.requests.isapi", iisStat.ISAPIExtensionRequestsPersec, tag))
			L = append(L, CounterValue("iis.errors.locked", iisStat.LockedErrorsPersec, tag))
			L = append(L, CounterValue("iis.requests.lock", iisStat.LockRequestsPersec, tag))
			L = append(L, CounterValue("iis.requests.mkcol", iisStat.MkcolRequestsPersec, tag))
			L = append(L, CounterValue("iis.requests.move", iisStat.MoveRequestsPersec, tag))
			L = append(L, CounterValue("iis.errors.notfound", iisStat.NotFoundErrorsPersec, tag))
			L = append(L, CounterValue("iis.requests.options", iisStat.OptionsRequestsPersec, tag))
			L = append(L, CounterValue("iis.requests.post", iisStat.PostRequestsPersec, tag))
			L = append(L, CounterValue("iis.requests.propfind", iisStat.PropfindRequestsPersec, tag))
			L = append(L, CounterValue("iis.requests.proppatch", iisStat.ProppatchRequestsPersec, tag))
			L = append(L, CounterValue("iis.requests.put", iisStat.PutRequestsPersec, tag))
			L = append(L, CounterValue("iis.requests.search", iisStat.SearchRequestsPersec, tag))
			L = append(L, CounterValue("iis.requests.trace", iisStat.TraceRequestsPersec, tag))
			L = append(L, CounterValue("iis.requests.unlock", iisStat.UnlockRequestsPersec, tag))
			L = append(L, GaugeValue("iis.service.uptime", iisStat.ServiceUptime, tag))
		}
	}
	return
}
