package g

import log "github.com/cihub/seelog"

func InitSeeLog() {
	//每个文件10MB大小，保存10份
	seelogConfig := `
	<seelog type="sync">
	    <outputs formatid="main">
		<console/>
		<rollingfile formatid="main" type="size" filename="./var/app.log" maxsize="10485760" maxrolls="10" />
	    </outputs>
	    <formats>
		<format id="main" format="%Date %Time [%LEV] %Msg%n"/>
	    </formats>
	</seelog>
	`
	logger, _ := log.LoggerFromConfigAsBytes([]byte(seelogConfig))
	log.ReplaceLogger(logger)
}
