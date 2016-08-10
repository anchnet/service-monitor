## mssql-Monitor-agent

监控sqlserver，采集相关信息上报给 open-falcon

version 信息上报给 smartAPI

以 agent 形式运行，提供简单的 http 信息维护接口

#### 上报字段
视版本和配置不同，采集到的 Metric 可能有所差别。

--------------------------------
| Counters | Type |Tag| Notes|
|-----|------|------|------|
|MsSQL.Lock_Waits/sec     |GUAGE|instance=instance|Lock_Waits/sec|
|MsSQL.Log_File(s)_Size_(KB)     |GUAGE|instance=instance|Log_File(s)_Size_(KB)|
|MsSQL.Log_File(s)_Used_Size_(KB)     |GUAGE|instance=instance|Log_File(s)_Used_Size_(KB)|
|MsSQL.Percent_Log_Used     |GUAGE|instance=instance|Log_File(s)_Used_Size_(KB)|
|MsSQL.Errors/sec     |GUAGE|error_type=error_type|Log_File(s)_Used_Size_(KB)|
|MsSQL.Batch_Requests/sec     |GUAGE|\|Batch_Requests/sec|
|MsSQL.Target_Server_Memory_(KB)     |GUAGE|\|Target_Server_Memory_(KB)|
|MsSQL.Total_Server_Memory_(KB)     |GUAGE|\|Total_Server_Memory_(KB)|
|MsSQL.IO_requests     |GUAGE|\|IO_requests|
|MsSQL.Connection     |GUAGE|\|Connections|
|MsSQL.Uptime    |GUAGE|\|Service Uptime|

其中Lock_Waits/sec …… Total_Server_Memory_(KB) 等通过查询sys.dm_os_performance_counters 获得，这需要服务器上开启性能计数器。

如果这部分指标缺失，请确认性能计数器是否正确开启。

#### 使用方式


配置文件请参照cfg.example.json，修改该文件名为cfg.json

```
{
	"debug": true,
	"logfile": "mssql.log",
	"hostname": "",
	"mssql":{
		"enabled": true,
		"addr":"127.0.0.1",
		"port":1433,
		"username":"sa",
		"password":"123456",
		"instance": [ //数据库实例名
	        "test"
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
cd %GOPATH%/src/github.com/51idc/service-monitor/mssql-monitor/
go get ./...
go build -o mssql-monitor.exe
将 mssql-monitor.exe cfg.json  打包，用于部署

```

#### 运行
以下命令需在管理员模式下运行开启 cmd/Powershell

先试运行一下
```
.\mssql-monitor.exe
2016/08/08 13:44:31 cfg.go:96: read config file: cfg.json successfully
2016/08/08 13:44:31 var.go:31: logging on iis.log
2016/08/08 13:44:31 http.go:64: listening :1990
```
等待1-2分钟，观察输出，确认运行正常
使用 [nssm](https://nssm.cc/) 注册为 Windows 服务。

```
.\nssm.exe install mssql-monitor
Administrator access is needed to install a service.
```
![](http://i.imgur.com/9hmkeOf.png)

