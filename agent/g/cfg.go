package g

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"os"
	"sync"

	"github.com/toolkits/file"
	"time"
)

type DiscoverConfig struct {
	Url   string `json:"url"`
	Cycle int    `json:"cycle"`
}

type PluginConfig struct {
	Enabled bool   `json:"enabled"`
	Dir     string `json:"dir"`
	Git     string `json:"git"`
	LogDir  string `json:"logs"`
}

type HeartbeatConfig struct {
	Enabled  bool   `json:"enabled"`
	Addr     string `json:"addr"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type TransferConfig struct {
	Enabled  bool     `json:"enabled"`
	Addrs    []string `json:"addrs"`
	Interval int      `json:"interval"`
	Timeout  int      `json:"timeout"`
}

type HttpConfig struct {
	Enabled  bool   `json:"enabled"`
	Listen   string `json:"listen"`
	Backdoor bool   `json:"backdoor"`
}

type CollectorConfig struct {
	IfacePrefix []string `json:"ifacePrefix"`
}

type GlobalConfig struct {
	Debug         bool             `json:"debug"`
	Hostname      string           `json:"hostname"`
	IP            string           `json:"ip"`
	Discover      *DiscoverConfig  `json:"discover"`
	Plugin        *PluginConfig    `json:"plugin"`
	Heartbeat     *HeartbeatConfig `json:"heartbeat"`
	Transfer      *TransferConfig  `json:"transfer"`
	Http          *HttpConfig      `json:"http"`
	SmartAPI      string           `json:"smartapi"`
	Collector     *CollectorConfig `json:"collector"`
	IgnoreMetrics map[string]bool  `json:"ignore"`
	Port          []string          `json:"port"`
	DialTimeout   time.Duration     `json:"dialTimeOut"`
	Process       map[string]bool  `json:"process"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func Hostname() (string, error) {
	hostname := Config().Hostname
	if hostname != "" {
		return hostname, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Info("ERROR: os.Hostname() fail", err)
	}
	return hostname, err
}

func IP() string {
	ip := Config().IP
	if ip != "" {
		// use ip in configuration
		return ip
	}

	if len(LocalIps) > 0 {
		ip = LocalIps[0]
	}

	return ip
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Error("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Error("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Error("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Error("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Info("read config file:", cfg, "successfully")
}
