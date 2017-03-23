package funcs

import (
	//	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

func mongo_serverStatus(Addr string, AuthDB string, Username string, Password string) (map[string]interface{}, error) {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{Addr},
		Timeout:  5 * time.Second,
		Database: AuthDB,
		Username: Username,
		Password: Password,
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	result := bson.M{}
	if err := session.DB("admin").Run(bson.D{{"serverStatus", 1}}, &result); err != nil {
		return nil, err
	} else {
		return result, err
	}
}
func mongo_replSetGetStatus(Addr string, AuthDB string, Username string, Password string) (map[string]interface{}, error) {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{Addr},
		Timeout:  60 * time.Second,
		Database: AuthDB,
		Username: Username,
		Password: Password,
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	result := bson.M{}
	if err := session.DB("admin").Run(bson.D{{"replSetGetStatus", 1}}, &result); err != nil {
		return nil, err
	} else {
		return result, err
	}
}
func mongo_version(serverStatus map[string]interface{}) string {
	if version, ok := serverStatus["version"]; ok {
		return version.(string)
	}
	return ""
}

func mongo_Metrics(serverStatus map[string]interface{}) (CounterMetrics map[string]int64, GaugeMetrics map[string]int64) {
	CounterMetrics = make(map[string]int64)
	GaugeMetrics = make(map[string]int64)
	if uptime, ok := serverStatus["uptime"]; ok {
		GaugeMetrics["uptime"] = int64(uptime.(float64))
	}
	if asserts, ok := serverStatus["asserts"]; ok {
		asserts_map := asserts.(bson.M)
		for asserts_metric, asserts_METRIC := range asserts_map {
			CounterMetrics["asserts_"+asserts_metric] = int64(asserts_METRIC.(int))
		}
	}
	if connections, ok := serverStatus["connections"]; ok {
		connections_map := connections.(bson.M)
		if current, ok := connections_map["current"]; ok {
			GaugeMetrics["connections_current"] = int64(current.(int))
		}
		if available, ok := connections_map["available"]; ok {
			GaugeMetrics["connections_available"] = int64(available.(int))
		}
		if totalCreated, ok := connections_map["totalCreated"]; ok {
			CounterMetrics["connections_totalCreated"] = formatInt(totalCreated)
		}
	}
	if extra_info, ok := serverStatus["extra_info"]; ok {
		extra_info_map := extra_info.(bson.M)
		if page_faults, ok := extra_info_map["page_faults"]; ok {
			CounterMetrics["page_faults"] = int64(page_faults.(int))
		}
	}
	if globalLock, ok := serverStatus["globalLock"]; ok {
		globalLock_map := globalLock.(bson.M)
		if currentQueue, ok := globalLock_map["currentQueue"]; ok {
			currentQueue_map := currentQueue.(bson.M)
			if total, ok := currentQueue_map["total"]; ok {
				GaugeMetrics["globalLock_currentQueue_total"] = int64(total.(int))
			}
			if readers, ok := currentQueue_map["readers"]; ok {
				GaugeMetrics["globalLock_currentQueue_readers"] = int64(readers.(int))
			}
			if writers, ok := currentQueue_map["writers"]; ok {
				GaugeMetrics["globalLock_currentQueue_writers"] = int64(writers.(int))
			}
		}
	}
	if locks, ok := serverStatus["locks"]; ok {
		locks_map := locks.(bson.M)
		for lock_type, LOCK_TYPE := range locks_map {
			LOCK_TYPE_map := LOCK_TYPE.(bson.M)
			for lock_metric, LOCK_METRIC := range LOCK_TYPE_map {
				LOCK_METRIC_map := LOCK_METRIC.(bson.M)
				if ISlock, ok := LOCK_METRIC_map["r"]; ok {
					CounterMetrics["locks_"+lock_type+"_"+lock_metric+"_ISlock"] = ISlock.(int64)
				}
				if IXlock, ok := LOCK_METRIC_map["w"]; ok {
					CounterMetrics["locks_"+lock_type+"_"+lock_metric+"_IXlock"] = IXlock.(int64)
				}
				if Slock, ok := LOCK_METRIC_map["R"]; ok {
					CounterMetrics["locks_"+lock_type+"_"+lock_metric+"_Slock"] = Slock.(int64)
				}
				if Xlock, ok := LOCK_METRIC_map["W"]; ok {
					CounterMetrics["locks_"+lock_type+"_"+lock_metric+"_Xlock"] = Xlock.(int64)
				}
			}
		}
	}
	if network, ok := serverStatus["network"]; ok {
		network_map := network.(bson.M)
		if bytesIn, ok := network_map["bytesIn"]; ok {
			switch bytesin := bytesIn.(type) {
			case int:
				CounterMetrics["network_bytesIn"] = int64(bytesin)
			case int64:
				CounterMetrics["network_bytesIn"] = bytesin
			}
		}
		if bytesOut, ok := network_map["bytesOut"]; ok {
			switch bytesout := bytesOut.(type) {
			case int:
				CounterMetrics["network_bytesOut"] = int64(bytesout)
			case int64:
				CounterMetrics["network_bytesOut"] = bytesout
			}
		}
		if numRequests, ok := network_map["numRequests"]; ok {
			switch numrequests := numRequests.(type) {
			case int:
				CounterMetrics["network_numRequests"] = int64(numrequests)
			case int64:
				CounterMetrics["network_numRequests"] = numrequests
			}
		}
	}
	if opcounters, ok := serverStatus["opcounters"]; ok {
		opcounters_map := opcounters.(bson.M)
		for opcounters_metric, opcounters_METRIC := range opcounters_map {
			CounterMetrics["opcounters_"+opcounters_metric] = int64(opcounters_METRIC.(int))
		}
	}
	if opcountersRepl, ok := serverStatus["opcountersRepl"]; ok {
		opcountersRepl_map := opcountersRepl.(bson.M)
		for opcountersRepl_metric, opcountersRepl_METRIC := range opcountersRepl_map {
			CounterMetrics["opcountersRepl_"+opcountersRepl_metric] = int64(opcountersRepl_METRIC.(int))
		}
	}
	if mem, ok := serverStatus["mem"]; ok {
		mem_map := mem.(bson.M)
		if resident, ok := mem_map["resident"]; ok {
			GaugeMetrics["mem_resident"] = int64(resident.(int) * 1024 * 1024)
		}
		if virtual, ok := mem_map["virtual"]; ok {
			GaugeMetrics["mem_virtual"] = int64(virtual.(int) * 1024 * 1024)
		}
		if mapped, ok := mem_map["mapped"]; ok {
			GaugeMetrics["mem_mapped"] = int64(mapped.(int) * 1024 * 1024)
		}
		if mappedWithJournal, ok := mem_map["mappedWithJournal"]; ok {
			GaugeMetrics["mem_mappedWithJournal"] = int64(mappedWithJournal.(int) * 1024 * 1024)
		}
	}
	if metrics, ok := serverStatus["metrics"]; ok {
		metrics_map := metrics.(bson.M)
		if cursor, ok := metrics_map["cursor"]; ok {
			cursor_map := cursor.(bson.M)
			if timeOut, ok := cursor_map["timeOut"]; ok {
				switch timeout := timeOut.(type) {
				case int:
					CounterMetrics["cursor_timedOut"] = int64(timeout)
				case int64:
					CounterMetrics["cursor_timedOut"] = timeout
				}
			}
			if open, ok := cursor_map["open"]; ok {
				open_map := open.(bson.M)
				if noTimeout, ok := open_map["noTimeout"]; ok {
					switch notimeout := noTimeout.(type) {
					case int:
						GaugeMetrics["cursor_open_noTimeout"] = int64(notimeout)
					case int64:
						GaugeMetrics["cursor_open_noTimeout"] = notimeout
					}

				}
				if total, ok := open_map["total"]; ok {
					switch open_total := total.(type) {
					case int:
						GaugeMetrics["cursor_open_total"] = int64(open_total)
					case int64:
						GaugeMetrics["cursor_open_total"] = open_total
					}
				}
				if pinned, ok := open_map["pinned"]; ok {
					switch open_pinned := pinned.(type) {
					case int:
						GaugeMetrics["cursor_open_pinned"] = int64(open_pinned)
					case int64:
						GaugeMetrics["cursor_open_pinned"] = open_pinned
					}
				}
			}
		}
	}
	if backgroundFlushing, ok := serverStatus["backgroundFlushing"]; ok {
		backgroundFlushing_map := backgroundFlushing.(bson.M)
		if flushes, ok := backgroundFlushing_map["flushes"]; ok {
			CounterMetrics["backgroundFlushing_flushes"] = int64(flushes.(int))
		}
		if last_ms, ok := backgroundFlushing_map["last_ms"]; ok {
			CounterMetrics["backgroundFlushing_last_ms"] = int64(last_ms.(int))
		}
		if average_ms, ok := backgroundFlushing_map["average_ms"]; ok {
			switch average_ms_time := average_ms.(type) {
			case int:
				GaugeMetrics["backgroundFlushing_average_ms"] = int64(average_ms_time)
			case float64:
				GaugeMetrics["backgroundFlushing_average_ms"] = int64(average_ms_time)
			}
		}
		if total_ms, ok := backgroundFlushing_map["total_ms"]; ok {
			GaugeMetrics["backgroundFlushing_total_ms"] = int64(total_ms.(int))
		}
	}
	if wiredTiger, ok := serverStatus["wiredTiger"]; ok {
		wiredTiger_map := wiredTiger.(bson.M)
		if cache, ok := wiredTiger_map["cache"]; ok {
			cache_map := cache.(bson.M)
			if wt_cache_used_total_bytes, ok := cache_map["bytes currently in the cache"]; ok {
				GaugeMetrics["wt_cache_used_total_bytes"] = int64(wt_cache_used_total_bytes.(int))
			}
			if wt_cache_dirty_bytes, ok := cache_map["tracked dirty bytes in the cache"]; ok {
				GaugeMetrics["wt_cache_dirty_bytes"] = int64(wt_cache_dirty_bytes.(int))
			}
			if wt_cache_readinto_bytes, ok := cache_map["bytes read into cache"]; ok {
				CounterMetrics["wt_cache_readinto_bytes"] = int64(wt_cache_readinto_bytes.(int))
			}
			if wt_cache_writtenfrom_bytes, ok := cache_map["bytes written from cache"]; ok {
				CounterMetrics["wt_cache_writtenfrom_bytes"] = int64(wt_cache_writtenfrom_bytes.(int))
			}
		}
		if concurrentTransactions, ok := wiredTiger_map["concurrentTransactions"]; ok {
			concurrentTransactions_map := concurrentTransactions.(bson.M)
			if write, ok := concurrentTransactions_map["write"]; ok {
				write_map := write.(bson.M)
				if available, ok := write_map["available"]; ok {
					GaugeMetrics["wt_concurrentTransactions_write"] = int64(available.(int))
				}
			}
			if read, ok := concurrentTransactions_map["read"]; ok {
				read_map := read.(bson.M)
				if available, ok := read_map["available"]; ok {
					GaugeMetrics["wt_concurrentTransactions_read"] = int64(available.(int))
				}
			}
		}
		if block_manager, ok := wiredTiger_map["block-manager"]; ok {
			block_manager_map := block_manager.(bson.M)
			if wt_bm_bytes_read, ok := block_manager_map["bytes read"]; ok {
				CounterMetrics["wt_bm_bytes_read"] = int64(wt_bm_bytes_read.(int))
			}
			if wt_bm_bytes_written, ok := block_manager_map["bytes written"]; ok {
				CounterMetrics["wt_bm_bytes_written"] = int64(wt_bm_bytes_written.(int))
			}
			if wt_bm_blocks_read, ok := block_manager_map["blocks read"]; ok {
				CounterMetrics["wt_bm_blocks_read"] = int64(wt_bm_blocks_read.(int))
			}
			if wt_bm_blocks_written, ok := block_manager_map["blocks written"]; ok {
				CounterMetrics["wt_bm_blocks_written"] = int64(wt_bm_blocks_written.(int))
			}
		}
	}
	return
}


func formatInt(value interface{}) int64{
    v := reflect.ValueOf(value)
    switch v.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16,
        reflect.Int32, reflect.Int64:
        return v.Int()
    default:
        return 0
    }
}