package g

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/toolkits/file"
)

type SmartAPIConfig struct {
	Enabled bool   `json:"enabled"`
	Url     string `json:"url"`
}

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type GlobalConfig struct {
	Debug    bool            `json:"debug"`
	SmartAPI *SmartAPIConfig `json:"smartAPI`
	Http     *HttpConfig     `json:"http"`
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
