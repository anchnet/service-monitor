package g

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/toolkits/file"
)

type QueueConfig struct {
	Sms  string `json:"sms"`
	Mail string `json:"mail"`
}

type RedisConfig struct {
	Addr    string `json:"addr"`
	MaxIdle int    `json:"maxIdle"`
}

type TrapServerConfig struct {
	TrapCommunity string `json:"trapcommunity`
	Community     string `json:"community"`
	Retry         int    `json:"retry`
	Timeout       int    `json:"timeout"`
	Listen        string `json:"listen"`
}

type UserConfig struct {
	Sms  []string `json:"sms"`
	Mail []string `json:"mail`
}

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type GlobalConfig struct {
	Debug      bool              `json:"debug"`
	Queue      *QueueConfig      `json:"queue"`
	Redis      *RedisConfig      `json:"redis"`
	User       *UserConfig       `json:"user"`
	TrapServer *TrapServerConfig `json:"trapServer"`
	Http       *HttpConfig       `json:"http"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
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
