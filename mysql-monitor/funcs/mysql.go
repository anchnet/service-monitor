package funcs

import (
	"fmt"
	"log"
	"time"

	"github.com/51idc/service-monitor/mysql-monitor/g"
	"github.com/open-falcon/common/model"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

func MysqlVersion(m *MysqlIns, username string, password string) (string, error) {
	db := mysql.New("tcp", "", fmt.Sprintf("%s:%d", m.Host, m.Port),
		username, password)
	db.SetTimeout(500 * time.Millisecond)
	if err := db.Connect(); err != nil {
		return "", err
	}
	defer db.Close()
	version, err := mysqlversion(m, db)
	return version, err
}

func MysqlStatus(m *MysqlIns, username string, password string) ([]*MetaData, error) {
	db := mysql.New("tcp", "", fmt.Sprintf("%s:%d", m.Host, m.Port),
		username, password)
	db.SetTimeout(500 * time.Millisecond)
	if err := db.Connect(); err != nil {
		return nil, err
	}
	defer db.Close()

	data := make([]*MetaData, 0)
	globalStatus, err := GlobalStatus(m, db)
	if err != nil {
		return nil, err
	}
	data = append(data, globalStatus...)

	globalVars, err := GlobalVariables(m, db)
	if err != nil {
		return nil, err
	}
	data = append(data, globalVars...)

	innodbState, err := innodbStatus(m, db)
	if err != nil {
		return nil, err
	}
	data = append(data, innodbState...)

	slaveState, err := slaveStatus(m, db)
	if err != nil {
		return nil, err
	}
	data = append(data, slaveState...)
	return data, err
}

func MysqlMetrics() (L []*model.MetricValue) {
	if !g.Config().Mysql.Enabled {
		log.Println("Mysql Monitor is disabled")
		return
	}
	Addr := g.Config().Mysql.Addr
	Port := g.Config().Mysql.Port
	Username := g.Config().Mysql.Username
	Password := g.Config().Mysql.Password

	debug := g.Config().Debug
	smartAPI_url := g.Config().SmartAPI.Url

	m := &MysqlIns{
		Host: Addr,
		Port: Port,
		Tag:  fmt.Sprintf("port=%d", Port),
	}

	if g.Config().SmartAPI.Enabled {
		endpoint, err := g.Hostname()
		version, err := MysqlVersion(m, Username, Password)
		if err == nil {
			smartAPI_Push(smartAPI_url, endpoint, version, debug)
		} else {
			log.Println(err)
		}
	}
	data, err := MysqlStatus(m, Username, Password)
	if err != nil {
		L = append(L, GaugeValue("mysql_alive_local", -1, m.Tag))
		if debug {
			log.Println("Mysql is not alive")
		}
		return
	}
	L = append(L, GaugeValue("mysql_alive_local", 1, m.Tag))
	if debug {
		log.Println("Mysql is alive")
	}
	for _, mstat := range data {
		if mstat.CounterType == "GAUGE" {
			L = append(L, GaugeValue(mstat.Metric, mstat.Value, m.Tag))
		}
		if mstat.CounterType == "COUNTER" {
			L = append(L, CounterValue(mstat.Metric, mstat.Value, m.Tag))
		}
	}
	return
}
