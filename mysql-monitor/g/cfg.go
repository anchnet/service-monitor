package g

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/toolkits/file"
)

type MysqlConfig struct {
	Enabled       bool            `json:"enabled"`
	Addr          string          `json:"addr"`
	Port          int             `json："port"`
	Username      string          `json:"username"`
	Password      string          `json:"password"`
	IgnoreMetrics map[string]bool `json:"ignore"`
}

type SmartAPIConfig struct {
	Enabled bool   `json:"enabled"`
	Url     string `json:"url"`
}

type TransferConfig struct {
	Enabled  bool   `json:"enabled"`
	Addrs    []string `json:"addrs"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type GlobalConfig struct {
	Debug    bool            `json:"debug"`
	Hostname string          `json:"hostname"`
	Mysql    *MysqlConfig    `json:"mysql"`
	SmartAPI *SmartAPIConfig `json:"smartAPI`
	Transfer *TransferConfig `json:"transfer"`
	Http     *HttpConfig     `json:"http"`
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
		log.Println("ERROR: os.Hostname() fail", err)
	}
	return hostname, err
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")

}
