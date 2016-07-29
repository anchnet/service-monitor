package funcs

import (
	"testing"
)

const (
	pushurl  = "http://127.0.0.1/api/sevice/version"
	Addr     = "127.0.0.1:6379"
	Password = ""
	DB       = 0
)

func Test_Redis_Info_Map(t *testing.T) {
	if info, err := GetRedisInfo(Addr, Password, DB); err != nil {
		t.Error(err)
	} else {
		t.Log("info:", info)
		Redis_Info := Redis_Info_Map(info)
		t.Log("Redis_Info:", Redis_Info)
	}
}

func Test_smartAPI_Push(t *testing.T) {
	endpoint := "test"
	version := "1.1.1"
	smartAPI_Push(pushurl, endpoint, version, true)

}
