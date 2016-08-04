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

```
iis-monitor.exe -c cfg.json
```

#### 验证

```
./iis-monitor.exe --check
```
