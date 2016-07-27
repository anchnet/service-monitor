## Tomcat-Monitor-agent

通过监控 Tomcat 自带的 statu 页面，采集相关信息上报给 open-falcon

version 信息上报给 smartAPI

以 agent 形式运行，提供简单的 http 信息维护接口

#### 上报字段
视版本不同，采集到的 Metric 可能有所差别。

--------------------------------
| key |  tag | type | note |
|-----|------|------|------|
|Tomcat.Monitor.alive|/|GUAGE|监控 agent 存活状态|
|Tomcat.Uptime|/|GUAGE|Tomcat Uptime|
|Tomcat.Jvm.Memory.Free|/|GAUGE|Jvm.Memory.Free|
|Tomcat.Jvm.Memory.Total|/|GAUGE|Jvm.Memory.Total|
|Tomcat.Jvm.Memory.Max|/|GAUGE|Jvm.Memory.Max|
|Tomcat.Jvm.Memory.usage|/|GAUGE|(Total-Free)/Total|
|Tomcat.Jvm.Memorypool.Initial|Name=Name,Type=Type|GAUGE|Jvm.Memorypool.Initial|
|Tomcat.Jvm.Memorypool.Committed|Name=Name,Type=Type|GAUGE|Jvm.Memorypool.Committed|
|Tomcat.Jvm.Memorypool.Max|Name=Name,Type=Type|GAUGE|Jvm.Jvm.Memorypool.Max|
|Tomcat.Jvm.Memorypool.Used|Name=Name,Type=Type|GAUGE|Jvm.Jvm.Memorypool.Used|
|Tomcat.Jvm.Memorypool.Usage|Name=Name,Type=Type|GAUGE|Usage/Max|
|Tomcat.Connector.ThreadInfo.MaxThreads|Connector=Connector|GAUGE|Connector.ThreadInfo.MaxThreads|
|Tomcat.Connector.ThreadInfo.CurrentThreadCount|Connector=Connector|GAUGE|Connector.ThreadInfo.CurrentThreadCount|
|Tomcat.Connector.ThreadInfo.CurrentThreadsBusy|Connector=Connector|GAUGE|Connector.ThreadInfo.CurrentThreadsBusy|
|Tomcat.Connector.RequestInfo.MaxTime|Connector=Connector|GAUGE|Connector.RequestInfo.MaxTime|
|Tomcat.Connector.RequestInfo.RequestCount|Connector=Connector|COUNTER|Connector.RequestInfo.RequestCount|
|Tomcat.Connector.RequestInfo.ErrorCount|Connector=Connector|COUNTER|Connector.RequestInfo.ErrorCount|
|Tomcat.Connector.RequestInfo.BytesReceived|Connector=Connector|COUNTER|Connector.RequestInfo.BytesReceived|
|Tomcat.Connector.RequestInfo.BytesSent|Connector=Connector|COUNTER|Connector.RequestInfo.BytesSent|


#### 使用方式
先配置 Tomcat， 为状态监控页增加一个用户
修改 `$TOMCAT_HOME/conf/tomcat_user.xml`
在 `<tomcat-users>` 标签下增加一个 `manger-gui` 用户
```
<tomcat-users>
……
<role rolename="manager-gui"/>
<user username="admin" password="manager" roles="manager-gui"/>
……
</tomcat-users>
```

配置文件请参照cfg.example.json，修改该文件名为cfg.json

```
{
	"debug": true,
	"hostname": "", #上报的 endpoint 名，如为空则是主机名
	"tomcat":{
		"enabled": true,
		"staturl": "http://127.0.0.1:8080/manager/status?XML=true",
		"username": "admin",
		"password": "manager"
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
cd $GOPATH/src/github.com/51idc/service-monitor/tomcat-monitor/
go get ./...
chmod +x control
./control build
./control pack
最后一步会pack出一个tar.gz的安装包，拿着这个包去部署服务即可

```

#### 进程管理

```
./control start 启动进程
./control stop 停止进程
./control restart 重启进程
./control status 查看进程状态
./control tail 用tail -f的方式查看var/app.log
```

#### 验证

```
./falcon-tomcat-monitor --check
```
