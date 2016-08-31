package trap

import (
	"log"
	"net"
	"os"
	"time"

	"strings"

	"github.com/51idc/service-monitor/snmptrap-server/g"
	"github.com/soniah/gosnmp"
)

func Start() {
	params := gosnmp.Default
	params.Logger = log.New(os.Stdout, "", 0)

	tl := gosnmp.TrapListener{
		OnNewTrap: myTrapHandler,
		Params:    params,
	}
	addr := g.Config().TrapServer.Listen
	if addr == "" {
		return
	}
	err := tl.Listen(addr)
	if err != nil {
		log.Panicf("error in listen: %s", err)
	}
}

func myTrapHandler(packet *gosnmp.SnmpPacket, addr *net.UDPAddr) {
	var err error
	var event Event
	event.IP = addr.IP.String()
	event.EventTime = time.Now().Unix()
	if packet.Community == g.Config().TrapServer.TrapCommunity {
		for _, v := range packet.Variables {
			if g.Config().Debug {
				log.Println(event.IP, v)
			}
			if v.Type == gosnmp.ObjectIdentifier {
				switch v.Value.(string) {
				case ".1.3.6.1.6.3.1.1.5.3":
					event.traptype = "linkDown"
				case ".1.3.6.1.6.3.1.1.5.4":
					event.traptype = "linkUp"
				default:
					return
				}
			}
			if strings.Contains(v.Name, ifindexoid) {
				ifindex := strings.Replace(v.Name, "."+ifindexoid, "", -1)
				ifname := ifnameoid + ifindex
				ifoperstatus := ifoperstatusoid + ifindex
				oids := []string{ifname, ifoperstatus}
				event.ifName, event.ifOperStatus, err = getIfinfo(oids, event.IP, g.Config().TrapServer.Community, g.Config().TrapServer.Timeout, g.Config().TrapServer.Retry)
			}
		}
		if err == nil {
			log.Println("put message in redis queue")
			sendevent(event, g.Config().User.Mail, g.Config().User.Sms)
		} else {
			log.Println(err)
		}
	} else {
		log.Println(event.IP, "trapcommunity not match: ", packet.Community)
	}
}
