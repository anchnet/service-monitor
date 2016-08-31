package trap

import (
	"fmt"
	"log"

	"github.com/51idc/service-monitor/snmptrap-server/g"
	"github.com/51idc/service-monitor/snmptrap-server/redis"
	"github.com/open-falcon/common/utils"
)

type Event struct {
	IP           string
	ifName       string
	ifOperStatus int
	traptype     string
	EventTime    int64
}

func (this *Event) FormattedTime() string {
	return utils.UnixTsFormat(this.EventTime)
}

func (this *Event) status() string {
	if this.ifOperStatus == 1 {
		return "OK"
	} else {
		return "PROBLEM"
	}
}

func BuildCommonSMSContent(event Event) string {
	return fmt.Sprintf(
		"[%s][%s][][%s %s %d ][%s]",
		event.status(),
		event.IP,
		event.ifName,
		event.traptype,
		event.ifOperStatus,
		event.FormattedTime(),
	)
}

func BuildCommonMailContent(event Event) string {
	return fmt.Sprintf(
		"%s\r\nIP: %s\r\nifName: %s\r\ntraptype: %s\r\nifOperStatus: %d\r\nTimestamp:%s\r\n",
		event.status(),
		event.IP,
		event.ifName,
		event.traptype,
		event.ifOperStatus,
		event.FormattedTime(),
	)
}

func GenerateSmsContent(event Event) string {
	return BuildCommonSMSContent(event)
}

func GenerateMailContent(event Event) string {
	return BuildCommonMailContent(event)
}

func sendevent(event Event, mail []string, sms []string) {

	smsContent := GenerateSmsContent(event)
	mailContent := GenerateMailContent(event)

	redis.WriteSms(sms, smsContent)
	redis.WriteMail(mail, smsContent, mailContent)
	if g.Config().Debug {
		log.Println(mailContent)
	}
}
