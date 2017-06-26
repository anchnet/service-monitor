package funcs

import (
	"database/sql"
	"strings"
	"time"

	"github.com/51idc/service-monitor/oracle-monitor/g"
	"golang.org/x/net/context"

	"github.com/open-falcon/common/model"
	_ "gopkg.in/rana/ora.v4"
)

type tablespace struct {
	TABLESPACE_NAME string
	USED_PERCENT    float64
}
type sysmetric struct {
	METRIC_NAME string
	VALUE       float64
}
type waitmetric struct {
	WAIT_CLASS           string
	AVERAGE_WAITER_COUNT float64
	DBTIME_IN_WAIT       float64
}

func oracleMetrics() (L []*model.MetricValue) {
	dsn := g.Config().Db.Dsn
	timeout := g.Config().Db.Timeout

	db, err := oracle_conn(dsn)
	if err != nil {
		g.Logger().Println(err)
		return
	}
	defer db.Close()

	debug := g.Config().Debug
	smartAPI_url := g.Config().SmartAPI.Url

	if g.Config().SmartAPI.Enabled {
		result, err := version_query(db, timeout)
		endpoint, _ := g.Hostname()
		if err == nil {
			version := result
			smartAPI_Push(smartAPI_url, endpoint, version, debug)
		} else {
			g.Logger().Println(err)
		}
	}
	database, err := database_query(db, timeout)
	instance, err := instance_query(db, timeout)
	database_tag := ""
	instance_tag := ""
	if err == nil {
		database_tag = "database=" + database
		instance_tag = "instance=" + instance
	} else {
		g.Logger().Println(err)
	}

	uptime, err := uptime_query(db, timeout)
	if err == nil {
		L = append(L, GaugeValue("Oracle.Uptime", uptime, database_tag, instance_tag))
	} else {
		g.Logger().Println(err)
	}
	tspaces, err := tablespace_query(db, timeout)
	if err == nil {
		for _, tspace := range tspaces {
			tspace_tag := "tablespace_name=" + tspace.TABLESPACE_NAME
			L = append(L, GaugeValue("Oracle.tablespace", tspace.USED_PERCENT, tspace_tag, database_tag, instance_tag))
		}
	} else {
		g.Logger().Println(err)
	}

	smetrics, err := sysmetric_query(db, timeout)
	if err == nil {
		for _, smetric := range smetrics {
			L = append(L, GaugeValue("Oracle.sysmetric."+smetric.METRIC_NAME, smetric.VALUE, database_tag, instance_tag))
		}
	} else {
		g.Logger().Println(err)
	}

	wmetrics, err := waitmetric_query(db, timeout)
	if err == nil {
		for _, wmetric := range wmetrics {
			wmetric_tag := "wait_class=" + wmetric.WAIT_CLASS
			L = append(L, GaugeValue("Oracle.waitmetric.avg_waiter_1m", wmetric.AVERAGE_WAITER_COUNT, wmetric_tag, database_tag, instance_tag))
			L = append(L, GaugeValue("Oracle.waitmetric.avg_dbtime_wait_1m", wmetric.DBTIME_IN_WAIT, wmetric_tag, database_tag, instance_tag))
		}
	} else {
		g.Logger().Println(err)
	}

	return
}

func formatMetric(metric string) string {
	metric = strings.Replace(metric, "<", "Below_", -1)
	metric = strings.Replace(metric, ">=", "Above_", -1)
	metric = strings.Replace(metric, "(", "", -1)
	metric = strings.Replace(metric, ")", "", -1)
	metric = strings.Replace(metric, "%", "Ratio", -1)
	metric = strings.TrimSpace(metric)
	metric = strings.Replace(metric, " ", "_", -1)
	return metric
}

func oracle_conn(dsn string) (*sql.DB, error) {
	db, err := sql.Open("ora", dsn)
	if err != nil {
		return nil, err
	}
	return db, err
}

func version_query(db *sql.DB, timeout int) (string, error) {
	sql := "select BANNER from v$version"
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var version string
	for rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			return "", err
		}
		break
	}
	return version, err
}

func uptime_query(db *sql.DB, timeout int) (float64, error) {
	sql := "select (sysdate - startup_time)*86400 uptime from sys.v_$instance"
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return 0.0, err
	}
	defer rows.Close()
	var uptime float64
	for rows.Next() {
		err = rows.Scan(&uptime)
		if err != nil {
			return 0.0, err
		}
		break
	}
	return uptime, err
}

func database_query(db *sql.DB, timeout int) (string, error) {
	sql := "select name from v$database"
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var database string
	for rows.Next() {
		err = rows.Scan(&database)
		if err != nil {
			return "", err
		}
		break
	}
	return database, err
}

func instance_query(db *sql.DB, timeout int) (string, error) {
	sql := "select instance_name from v$instance"
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var instance string
	for rows.Next() {
		err = rows.Scan(&instance)
		if err != nil {
			return "", err
		}
		break
	}
	return instance, err
}

func tablespace_query(db *sql.DB, timeout int) ([]tablespace, error) {
	sql := "select TABLESPACE_NAME , USED_PERCENT from dba_tablespace_usage_metrics"
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	tspaces := []tablespace{}
	var tspace tablespace
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&tspace.TABLESPACE_NAME, &tspace.USED_PERCENT)
		if err != nil {
			return nil, err
		}
		tspace.TABLESPACE_NAME = formatMetric(tspace.TABLESPACE_NAME)
		tspaces = append(tspaces, tspace)
	}
	return tspaces, err
}

func sysmetric_query(db *sql.DB, timeout int) ([]sysmetric, error) {
	sql := "select METRIC_NAME , VALUE from v$sysmetric"
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	smetrics := []sysmetric{}
	var smetric sysmetric
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&smetric.METRIC_NAME, &smetric.VALUE)
		if err != nil {
			return nil, err
		}
		smetric.METRIC_NAME = formatMetric(smetric.METRIC_NAME)
		smetrics = append(smetrics, smetric)
	}
	return smetrics, err
}

func waitmetric_query(db *sql.DB, timeout int) ([]waitmetric, error) {
	sql := "select b.WAIT_CLASS , a.AVERAGE_WAITER_COUNT , a.DBTIME_IN_WAIT  from v$waitclassmetric a, v$system_wait_class b where a.WAIT_CLASS_ID = b.WAIT_CLASS_ID and WAIT_CLASS <> 'Idle'"
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	wmetrics := []waitmetric{}
	var wmetric waitmetric
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&wmetric.WAIT_CLASS, &wmetric.AVERAGE_WAITER_COUNT, &wmetric.DBTIME_IN_WAIT)
		if err != nil {
			return nil, err
		}
		wmetric.WAIT_CLASS = formatMetric(wmetric.WAIT_CLASS)
		wmetrics = append(wmetrics, wmetric)
	}
	return wmetrics, err
}
