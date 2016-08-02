package funcs

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/51idc/service-monitor/redis-monitor/g"
	"github.com/open-falcon/common/model"
	"gopkg.in/redis.v4"
)

func GetRedisInfo(Addr string, Password string, DB int) (string, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password, // no password set
		DB:       DB,       // use default DB
	})
	err := client.Ping().Err()
	info := client.Info()
	return fmt.Sprintln(info), err
}

func Redis_Info_Map(info string) map[string]string {
	Redis_Info := make(map[string]string)
	Info := strings.TrimLeft(info, "info: ")
	str := strings.Split(Info, "\n")
	for _, line := range str {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "#") {
			continue
		}
		if !strings.Contains(line, ":") {
			continue
		}
		v := strings.Split(line, ":")
		Redis_Info[v[0]] = v[1]
	}
	return Redis_Info
}

func RedisMetrics() (L []*model.MetricValue) {
	if !g.Config().Redis.Enabled {
		log.Println("Redis Monitor is disabled")
		return
	}
	Addr := g.Config().Redis.Addr
	Password := g.Config().Redis.Password
	DB := g.Config().Redis.Db

	info, err := GetRedisInfo(Addr, Password, DB)
	if err != nil {
		log.Println("Redis Connect Error: ", err)
		return
	}
	Port := strings.Split(Addr, ":")[1]
	Redis_Info := Redis_Info_Map(info)
	debug := g.Config().Debug
	smartAPI_url := g.Config().SmartAPI.Url

	if g.Config().SmartAPI.Enabled {
		version := Redis_Info["redis_version"]
		endpoint, err := g.Hostname()
		if err == nil {
			smartAPI_Push(smartAPI_url, endpoint, version, debug)
		} else {
			log.Println(err)
		}
	}
	for index, value := range Redis_Info {
		if Type, ok := Mertics[index]; ok {
			if index == "uptime_in_seconds" {
				value, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Println(index, err)
					continue
				}
				L = append(L, GaugeValue("Redis.Uptime", value, "Port="+Port))
				continue
			}
			switch Type {
			case "GUAGE":
				value, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Println(index, err)
					continue
				}
				L = append(L, GaugeValue("Redis."+index, value, "Port="+Port))
				if index == "maxmemory" {
					if value > 0 {
						maxmemory := value
						used_memory, err := strconv.ParseFloat(Redis_Info["used_memory"], 64)
						if err != nil {
							log.Println("used_memory", err)
						}
						used_memory_per := used_memory / maxmemory
						used_memory_pct := int(used_memory_per * 100)
						L = append(L, GaugeValue("Redis.used_memory_pct", used_memory_pct, "Port="+Port))
					}
				}
			case "COUNTER":
				value, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Println(index, err)
					continue
				}
				L = append(L, CounterValue("Redis."+index, value, "Port="+Port))
			}
		}
	}
	keyspace_hits, err := strconv.ParseFloat(Redis_Info["keyspace_hits"], 64)
	keyspace_misses, err := strconv.ParseFloat(Redis_Info["keyspace_misses"], 64)
	if err == nil && (keyspace_hits+keyspace_misses) != 0 {
		keyspace_hit_rat := keyspace_hits / (keyspace_hits + keyspace_misses)
		keyspace_hit_ratio := int(keyspace_hit_rat * 100)
		L = append(L, GaugeValue("Redis.keyspace_hit_ratio", keyspace_hit_ratio, "Port="+Port))
	} else {
		L = append(L, GaugeValue("Redis.keyspace_hit_ratio", 0, "Port="+Port))
	}

	return
}
