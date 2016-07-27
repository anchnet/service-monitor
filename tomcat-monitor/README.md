## Nginx-Monitor-agent

通过监控 Nginx 自带的 statu 页面，采集相关信息上报给 open-falcon

version 信息上报给 smartAPI

以 agent 形式运行，提供简单的 http 信息维护接口

#### 上报字段
--------------------------------
| key |  tag | type | note |
|-----|------|------|------|
|Nginx.Monitor.alive|/|GUAGE|监控 agent 的存活状态|
|Nginx.Uptime|/|GUAGE|运行时长|
|Nginx.ActiveConn|/|GAUGE|当前活跃连接数|
|Nginx.ServerAccepts|/|COUNTER|请求成功数|
|Nginx.ServerHandled|/|COUNTER|请求挂起数|
|Nginx.ServerRequests|/|COUNTER|请求数 |
|Nginx.ServerReading|/|GAUGE|读的 worker 数|
|Nginx.ServerWaiting|/|GAUGE|等待请求的 worker 数|
|Nginx.ServerWriting|/|GAUGE|写的 worker 数|


#### 使用方式
先配置 Nginx ，开启状态监控页
修改 nginx.conf 增加如下配置
```
      location /status {
              stub_status on;
              access_log /var/log/nginx/status.log;
              auth_basic "NginxStatus";
      }
```


配置文件请参照cfg.example.json，修改该文件名为cfg.json

```
{
	"debug": true,
	"hostname": "", #上报的 endpoint 名，如为空则是主机名
	"nginx":{
		"enabled": true,
		"staturl": "http://127.0.0.1/status",
		"pid": "/var/run/nginx.pid"
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
cd $GOPATH/src/github.com/51idc/service-monitor/nginx-monitor/
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
./falcon-nginx-monitor --check
```
