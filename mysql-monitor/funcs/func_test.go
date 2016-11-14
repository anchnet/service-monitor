package funcs

import (
	"fmt"
	"testing"
)

const (
	pushurl  = "http://127.0.0.1/api/sevice/version"
	Addr     = "192.168.11.136"
	Port     = 3306
	Username = "root"
	Password = ""
)

func Test_mysql_stat(t *testing.T) {
	m := &MysqlIns{
		Host: Addr,
		Port: Port,
		Tag:  fmt.Sprintf("port=%d", Port),
	}
	data, err := MysqlStatus(m, Username, Password)
	t.Log(data)
	t.Error(err)
	version, err := MysqlVersion(m, Username, Password)
	t.Log(version)
	t.Error(err)

}

func Test_smartAPI_Push(t *testing.T) {
	endpoint := ""
	version := "1.1.1"
	smartAPI_Push(pushurl, endpoint, version, true)

}
