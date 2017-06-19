## oracle-Monitor-agent

监控oracle，采集相关信息上报给 open-falcon

version 信息上报给 smartAPI

以 agent 形式运行，提供简单的 http 信息维护接口

#### 上报字段
视版本和配置不同，采集到的 Metric 可能有所差别。

--------------------------------
| Counters | Type |Tag| Notes|
|-----|------|------|------|



#### 使用方式


配置文件请参照cfg.example.json，修改该文件名为cfg.json

```
{
  "debug": true,
  "logfile": "oracle.log",
  "hostname": "",
  "db": {
   	"dsn": "c##test/test@127.0.0.1:1521/orcl",
    "timeout": 5
   },
  "smartAPI": {
    "enabled": true,
    "url": "http://127.0.0.1/api/service/version"
  },
  "transfer": {
    "enabled": true,
    "addrs": [
      "127.0.0.1:8433",
      "127.0.0.1:8433"
    ],
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
todo
```

#### 运行
以下命令需在管理员模式下运行开启 cmd/Powershell

先试运行一下
```
PS C:\Users\Administrator\Desktop\monitor> .\oracle-monitor.exe
2017/06/20 00:21:37 cfg.go:96: read config file: cfg.json successfully
2017/06/20 00:21:37 var.go:24: logging on oracle.log
2017/06/20 00:21:37 http.go:64: listening :1990
```
等待1-2分钟，观察输出，确认运行正常
使用 [nssm](https://nssm.cc/) 注册为 Windows 服务。

```
.\nssm.exe install oracle-monitor
Administrator access is needed to install a service.
```

