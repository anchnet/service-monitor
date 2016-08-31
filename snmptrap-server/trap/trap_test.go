package trap

import (
	"testing"
	"time"
)

const (
	ip        = "192.168.31.101"
	community = "ecnu-changpan"
)

func Test_runsnmp(t *testing.T) {
	oids := []string{"1.3.6.1.2.1.2.2.1.1.1", "1.3.6.1.2.1.2.2.1.2.1"}
	ifname, ifoperstatus, err := getIfinfo(oids, ip, community, 2, 3)
	t.Log(ifname)
	t.Log(ifoperstatus)
	t.Error(err)
}
func Test_GenerateSmsContent(t *testing.T) {
	var event Event
	event.IP = "192.168.31.101"
	event.ifName = "Fa0/1"
	event.traptype = "linkUp"
	event.ifOperStatus = 1
	event.EventTime = time.Now().Unix()
	s := GenerateSmsContent(event)
	t.Log(s)
}

func Test_GenerateMailContent(t *testing.T) {
	var event Event
	event.IP = "192.168.31.101"
	event.ifName = "Fa0/1"
	event.traptype = "linkUp"
	event.ifOperStatus = 1
	event.EventTime = time.Now().Unix()
	s := GenerateMailContent(event)
	t.Log(s)
}
