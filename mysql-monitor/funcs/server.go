package funcs

import (
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

func GlobalStatus(m *MysqlIns, db mysql.Conn) ([]*MetaData, error) {
	return mysqlState(m, db, "SHOW /*!50001 GLOBAL */ STATUS")
}

func GlobalVariables(m *MysqlIns, db mysql.Conn) ([]*MetaData, error) {
	return mysqlState(m, db, "SHOW /*!50001 GLOBAL */ VARIABLES")
}
func mysqlversion(m *MysqlIns, db mysql.Conn) (string, error) {
	rows, _, err := db.Query("select version()")
	if err != nil {
		return "", err
	}
	version := rows[0].Str(0)
	return version, err

}
func mysqlState(m *MysqlIns, db mysql.Conn, sql string) ([]*MetaData, error) {
	rows, _, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	data := make([]*MetaData, len(rows))
	i := 0
	for _, row := range rows {
		key_ := row.Str(0)
		v, err := row.Int64Err(1)
		// Ignore non digital value
		if err != nil {
			continue
		}

		data[i] = NewMetric(key_)
		data[i].SetValue(v)
		i++
	}
	return data[:i], nil
}
