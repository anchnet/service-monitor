package trap

import (
	"time"

	"github.com/soniah/gosnmp"
)

func runsnmp(oids []string, ip string, community string, timeout int, retry int) (snmpPDU []gosnmp.SnmpPDU, err error) {
	snmpPDU = []gosnmp.SnmpPDU{}
	params := &gosnmp.GoSNMP{
		Target:    ip,
		Port:      161,
		Community: community,
		Version:   gosnmp.Version2c,
		Retries:   retry,
		Timeout:   time.Duration(timeout) * time.Second,
	}
	err = params.Connect()
	if err != nil {
		return
	}
	defer params.Conn.Close()

	result, err := params.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		return
	}
	snmpPDU = result.Variables

	return
}
func getIfinfo(oids []string, ip string, community string, timeout int, retry int) (ifname string, ifoperstatus int, err error) {
	var snmpPDU []gosnmp.SnmpPDU
	snmpPDU, err = runsnmp(oids, ip, community, timeout, retry)

	if err == nil {
		for _, variable := range snmpPDU {
			switch variable.Type {
			case gosnmp.OctetString:
				ifname = string(variable.Value.([]byte))
			default:
				ifoperstatus = variable.Value.(int)
			}
		}
	}
	return

}
