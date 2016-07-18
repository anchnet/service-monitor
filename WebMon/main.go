package main

import (
	"flag"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	goconf "github.com/akrennmair/goconf"
)

type Cfg struct {
	LogFile        string
	LogLevel       int
	PushUrl        string
	Endpoint       string
	NginxStatUrl   string
	NginxEnabled   int
	ApacheStatUrl  string
	ApacheEnabled  int
	TomcatEnabled  int
	TomcatStatUrl  string
	TomcatUsername string
	TomcatPassword string
	Interval       int
}

var cfg Cfg

func init() {
	var cfgFile string
	flag.StringVar(&cfgFile, "c", "WebMon.cfg", "WebMon configure file")
	flag.Parse()

	if _, err := os.Stat(cfgFile); err != nil {
		if os.IsNotExist(err) {
			log.WithField("cfg", cfgFile).Fatalf("WebMon config file does not exists: %v", err)
		}
	}

	if err := cfg.readConf(cfgFile); err != nil {
		log.Fatalf("Read configure file failed: %v", err)
	}

	// Init log file
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.Level(cfg.LogLevel))

	if cfg.LogFile != "" {
		f, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			log.SetOutput(f)
			return
		}
	}
	log.SetOutput(os.Stderr)
}

func (conf *Cfg) readConf(file string) error {
	c, err := goconf.ReadConfigFile(file)
	if err != nil {
		return err
	}

	conf.LogFile, err = c.GetString("default", "log_file")
	if err != nil {
		return err
	}

	conf.LogLevel, err = c.GetInt("default", "log_level")
	if err != nil {
		return err
	}

	conf.PushUrl, err = c.GetString("default", "pushurl")
	if err != nil {
		return err
	}

	conf.Endpoint, err = c.GetString("default", "endpoint")
	if err != nil {
		return err
	}

	conf.Interval, err = c.GetInt("default", "interval")
	if err != nil {
		return err
	}

	conf.NginxEnabled, err = c.GetInt("nginx", "enabled")
	if err != nil {
		return err
	}

	conf.NginxStatUrl, err = c.GetString("nginx", "staturl")
	if err != nil {
		return err
	}

	conf.ApacheEnabled, err = c.GetInt("apache", "enabled")
	if err != nil {
		return err
	}

	conf.ApacheStatUrl, err = c.GetString("apache", "staturl")
	if err != nil {
		return err
	}
	conf.TomcatEnabled, err = c.GetInt("tomcat", "enabled")
	if err != nil {
		return err
	}

	conf.TomcatStatUrl, err = c.GetString("tomcat", "staturl")
	if err != nil {
		return err
	}

	conf.TomcatUsername, err = c.GetString("tomcat", "username")
	if err != nil {
		return err
	}

	conf.TomcatPassword, err = c.GetString("tomcat", "password")
	if err != nil {
		return err
	}

	return err
}

func timeout() {
	time.AfterFunc(TIME_OUT*time.Second, func() {
		log.Errorf("Execute timeout")
		os.Exit(1)
	})
}

func NginxAlive(url string, ok bool) {
	data := NewMetric("Nginx.alive")
	if ok {
		data.SetValue(1)
	} else {
		data.SetValue(-1)
	}
	_, err := sendData([]*MetaData{data})
	if err != nil {
		log.Errorf("Send alive data failed: %v", err)
		return
	}
	log.Infof("Alive data response Nginx")
}

func TomcatAlive(url string, ok bool) {
	data := NewMetric("Tomcat.alive")
	if ok {
		data.SetValue(1)
	} else {
		data.SetValue(-1)
	}
	_, err := sendData([]*MetaData{data})
	if err != nil {
		log.Errorf("Send alive data failed: %v", err)
		return
	}
	log.Infof("Alive data response Tomcat")
}

func ApacheAlive(url string, ok bool) {
	data := NewMetric("Apache.alive")
	if ok {
		data.SetValue(1)
	} else {
		data.SetValue(-1)
	}
	_, err := sendData([]*MetaData{data})
	if err != nil {
		log.Errorf("Send alive data failed: %v", err)
		return
	}
	log.Infof("Alive data response Apache")
}

func FetchNginxData(url string) (err error) {
	defer func() {
		NginxAlive(url, err == nil)
	}()

	respbody, resp_code, err := httpGet(url)
	if err != nil {
		log.Error(err)
		return
	}
	if resp_code != 200 {
		log.Errorf("Http Statu Page Open Error")
		return
	}
	stat, err := nginx_status(respbody)
	if err != nil {
		log.Error(err)
		return
	}
	data := nginx_data(stat)

	_, err = sendData(data)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Send response")
	return
}

func FetchApacheData(url string) (err error) {
	defer func() {
		ApacheAlive(url, err == nil)
	}()

	respbody, resp_code, err := httpGet(url)
	if err != nil {
		log.Error(err)
		return
	}
	if resp_code != 200 {
		log.Errorf("Http Statu Page Open Error")
		return
	}
	stat, err := apache_status(respbody)
	if err != nil {
		log.Error(err)
		return
	}
	data := apache_data(stat)

	_, err = sendData(data)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Send response")
	return
}

func FetchTomcatData(username string, password string, url string) (err error) {
	defer func() {
		TomcatAlive(url, err == nil)
	}()

	respbody, resp_code, err := TomcathttpGet(username, password, url)
	if err != nil {
		log.Error(err)
		return
	}
	if resp_code != 200 {
		log.Errorf("Http Statu Page Open Error")
		return
	}
	stat, err := xml_struct(respbody)
	if err != nil {
		log.Error(err)
		return
	}
	data := tomcat_data(stat)

	_, err = sendData(data)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Send response")
	return
}

func main() {
	log.Info("Web Monitor for falcon")
	go timeout()
	if cfg.NginxEnabled == 1 {
		err := FetchNginxData(cfg.NginxStatUrl)
		if err != nil {
			log.Error(err)
		}
	}
	if cfg.ApacheEnabled == 1 {
		err := FetchApacheData(cfg.ApacheStatUrl)
		if err != nil {
			log.Error(err)
		}
	}
	if cfg.TomcatEnabled == 1 {
		err := FetchTomcatData(cfg.TomcatUsername, cfg.TomcatPassword, cfg.TomcatStatUrl)
		if err != nil {
			log.Error(err)
		}
	}
}
