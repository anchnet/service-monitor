## Apache-Monitor-agent

通过监控 Apache 自带的 statu 页面，采集相关信息上报给 open-falcon

version 信息上报给 smartAPI

以 agent 形式运行，提供简单的 http 信息维护接口

#### 上报字段
视版本不同，采集到的 Metric 可能有所差别。

--------------------------------
| key |  tag | type | note |
|-----|------|------|------|
|Apache.Monitor.alive|/|GUAGE|监控 agent 存活状态|
|Apache.Total_Access|/|GAUGE|总访问数|
|Apache.Total_kBytes|/|GUAGE|总数据量，单位为 kb|
|Apache.CPULoad|/|GUAGE|CPU 负载|
|APache.Uptime|/|GUAGE|运行时长|
|Apache.ReqPerSec|/|GAUGE|每秒请求数|
|Apache.BytesPerSec|/|GAUGE|每秒数据量|
|Apache.BytesPerReq|/|GAUGE|每次请求数据量|
|Apache.BusyWorkers|/|GUAGE|在忙的 Worker 数|
|Apache.IdleWorkers|/|GUAGE|空闲的 Worker 数|
|Apache.ConnsTotal|/|GUAGE|当前连接数|
|Apache.ConnsAsyncWriting|/|GUAGE|异步写连接数|
|Apache.ConnsAsyncKeepAlive|/|GUAGE|异步 keepAlive 连接数|
|Apache.ConnsAsyncClosing|/|GUAGE|异步关闭中的连接数|
|Apache.Waiting_for_Connection|/|GUAGE|等待连接的 Worker 数|
|Apache.Starting_up|/|GUAGE|启动中的 Worker 数|
|Apache.Reading_Request|/|GUAGE|读请求中的 Worker 数|
|Apache.Sending_Reply|/|GUAGE|发送响应中的 Worker 数|
|Apache.Keepalive_read|/|GUAGE|keepalive 的 Worker 数|
|Apache.DNS_Lookup|/|GUAGE|DNS 解析中的 Worker 数|
|Apache.Closing_connection|/|GUAGE|关闭连接中的 Worker 数|
|Apache.Logging|/|GUAGE|写日志中的 Worker 数|
|Apache.Gracefully_finishing|/|GUAGE|Gracefully finishing 的 Worker 数|
|Apache.Idle_cleanup_of_worker|/|GUAGE|Idle cleanup of worker|
|Apache.Open_slot_with_no_current_process|/|GUAGE|Open slot with no current process|



#### 使用方式
先配置 Apache ，开启状态监控页
增加如下配置

```
<location /server-status>
         SetHandler server-status
         Order Deny,Allow  
         Deny from all
         Allow from 127.0.0.1
</location>
ExtendedStatus On
```


配置文件请参照cfg.example.json，修改该文件名为cfg.json

```
{
	"debug": true,
	"hostname": "", #上报的 endpoint 名，如为空则是主机名
	"apache":{
		"enabled": true,
		"staturl": "http://127.0.0.1/status",
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
cd $GOPATH/src/github.com/51idc/service-monitor/apache-monitor/
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
./falcon-apache-monitor --check
```
