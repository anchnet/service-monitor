package funcs

import (
	"testing"
)

const (
	pushurl  = "http://127.0.0.1/api/sevice/version"
	server   = "127.0.0.1"
	port     = 1433
	user     = "sa"
	password = "123456"
	encrypt  = "disable"
)

func Test_in_array(t *testing.T) {
	instance := []string{"test"}
	re := in_array("test", instance)
	t.Log(re)
}

func Test_performance_query(t *testing.T) {
	instance := []string{"_Total", "test"}
	db, err := mssql_conn(server, port, user, password, encrypt)
	if err != nil {
		t.Error(err)
	}
	result, err := performance_query(db, instance)
	t.Log(result)
	t.Error(err)
	db.Close()
}

func Test_io_req_query(t *testing.T) {
	db, err := mssql_conn(server, port, user, password, encrypt)
	if err != nil {
		t.Error(err)
	}
	result, err := io_req_query(db)
	t.Log(result)
	t.Error(err)
	db.Close()
}

func Test_conn_query(t *testing.T) {
	db, err := mssql_conn(server, port, user, password, encrypt)
	if err != nil {
		t.Error(err)
	}
	result, err := conn_query(db)
	t.Log(result)
	t.Error(err)
	db.Close()
}

func Test_version_query(t *testing.T) {
	db, err := mssql_conn(server, port, user, password, encrypt)
	if err != nil {
		t.Error(err)
	}
	result, err := version_query(db)
	t.Log(result)
	t.Error(err)
	db.Close()
}

func Test_uptime_query(t *testing.T) {
	db, err := mssql_conn(server, port, user, password, encrypt)
	if err != nil {
		t.Error(err)
	}
	result, err := uptime_query(db)
	t.Log(result)
	t.Error(err)
	db.Close()
}

func Test_smartAPI_Push(t *testing.T) {
	endpoint := "test"
	version := "1.1.1"
	smartAPI_Push(pushurl, endpoint, version, true)
}
