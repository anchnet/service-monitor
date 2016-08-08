## iis-Monitor-agent

通过监控 iis info 命令，采集相关信息上报给 open-falcon

version 信息上报给 smartAPI

以 agent 形式运行，提供简单的 http 信息维护接口

#### 上报字段
视版本和配置不同，采集到的 Metric 可能有所差别。

--------------------------------
| Counters | Type |Tag| Notes|
|-----|------|------|------|
|IIs.Total_Bytes_Received  |COUNTER|site=site|Bytes Received/sec |
|IIs.Total_Bytes_Sent      |COUNTER|site=site|Total Bytes Sent/sec|
|IIs.Total_Delete_Requests      |COUNTER|site=site|Delete Requests/sec|
|IIs.Total_Get_Requests     |COUNTER|site=site|Get Requests/sec|
|IIs.Total_Post_Requests     |COUNTER|site=site|Post Requests/sec|
|IIs.Total_Put_Requests     |COUNTER|site=site|Put Requests/sec|
|IIs.Total_Not_Found_Errors     |COUNTER|site=site|Not Found Errors/sec|
|IIs.Maximum_Connections     |GUAGE|site=site|Maximum Connections|
|IIs.Current_Connections     |GUAGE|site=site|Current Connections|
|IIs.Uptime     |GUAGE|site=site|Service Uptime|


#### 使用方式


配置文件请参照cfg.example.json，修改该文件名为cfg.json

```
{
	"debug": true,
	"logfile": "iis.log",
	"hostname": "",
	"iis":{
		"enabled": true,
		"websites": [
	        "Default Web Site" //iis 内站点的列表，可以监控某个具体站点的数据
	    ]
 	}, 
	"smartAPI": {
		"enabled": true,
		"url": "http://127.0.0.1/api/sevice/version"
	},
    "transfer": {
        "enabled": true,
        "addr": "127.0.0.1:8433",
        "interval": 60,
        "timeout": 1000
    },
    "http": {
        "enabled": true,
        "listen": ":1990"
    }
}



```

#### http 信息维护接口

```
curl http://127.0.0.1:1990/health
正常则返回 ok

curl http://127.0.0.1:1990/version
返回版本

curl http://127.0.0.1:1990/workdir
返回工作目录
 
curl http://127.0.0.1:1990/config
返回配置
```

#### 源码安装

```
cd %GOPATH%/src/github.com/51idc/service-monitor/iis-monitor/
go get ./...
go build -o iis-monitor.exe
将 iis-monitor.exe cfg.json scrips/ 打包，用于部署

```

#### 运行
以下命令需在管理员模式下运行开启命令行/Powershell

先试运行一下
```
.\iis-monitor.exe
2016/08/08 13:44:31 cfg.go:96: read config file: cfg.json successfully
2016/08/08 13:44:31 var.go:31: logging on iis.log
2016/08/08 13:44:31 http.go:64: listening :1990
```
等待1-2分钟，观察输出，确认运行正常
使用 [nssm](https://nssm.cc/) 注册为 Windows 服务。

```
.\nssm.exe install iis-monitor
Administrator access is needed to install a service.
```
![](http://i.imgur.com/9hmkeOf.png)

