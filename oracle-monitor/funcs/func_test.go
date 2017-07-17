package funcs

import (
	"testing"
)

const (
	dsn = "system/test123@127.0.0.1:1521/orcl"
)

func Test_version_query(t *testing.T) {
	db, err := oracle_conn(dsn)
	if err != nil {
		t.Error(err)
	}
	result, err := version_query(db, 5)
	t.Log(result)
	t.Error(err)
	db.Close()
}

func Test_database_query(t *testing.T) {
	db, err := oracle_conn(dsn)
	if err != nil {
		t.Error(err)
	}
	result, err := database_query(db, 5)
	t.Log(result)
	t.Error(err)
	db.Close()
}

func Test_instance_query(t *testing.T) {
	db, err := oracle_conn(dsn)
	if err != nil {
		t.Error(err)
	}
	result, err := instance_query(db, 5)
	t.Log(result)
	t.Error(err)
	db.Close()
}

func Test_uptime_query(t *testing.T) {
	db, err := oracle_conn(dsn)
	if err != nil {
		t.Error(err)
	}
	result, err := uptime_query(db, 5)
	t.Log(result)
	t.Error(err)
	db.Close()
}

func Test_tablespace_query(t *testing.T) {
	db, err := oracle_conn(dsn)
	if err != nil {
		t.Error(err)
	}
	result, err := tablespace_query(db, 5)
	t.Log(result)
	t.Error(err)
	db.Close()
}

func Test_sysmetric_query(t *testing.T) {
	db, err := oracle_conn(dsn)
	if err != nil {
		t.Error(err)
	}
	result, err := sysmetric_query(db, 5)
	t.Log(result)
	t.Error(err)
	db.Close()
}

func Test_waitmetric_query(t *testing.T) {
	db, err := oracle_conn(dsn)
	if err != nil {
		t.Error(err)
	}
	result, err := waitmetric_query(db, 5)
	t.Log(result)
	t.Error(err)
	db.Close()
}
