## Web-Monitor

通过监控各 web 应用 自带的 statu 页面，采集相关信息上报给 open-falcon

#### 上报字段
#### Nginx
--------------------------------
| key |  tag | type | note |
|-----|------|------|------|
|Nginx.ActiveConn|/|GAUGE|当前活跃连接数|
|Nginx.ServerAccepts|/|COUNTER|请求成功数|
|Nginx.ServerHandled|/|COUNTER|请求挂起数|
|Nginx.ServerRequests|/|COUNTER|请求数 |
|Nginx.ServerReading|/|GAUGE|读请求数|
|Nginx.ServerWaiting|/|GAUGE|等待请求数|
|Nginx.ServerWriting|/|GAUGE|写请求数|


#### 使用方式
#### Nginx
先配置 Nginx ，开启状态监控页
修改 nginx.conf 增加如下配置
```
      location /status {
              stub_status on;
              access_log /var/log/nginx/status.log;
              auth_basic "NginxStatus";
      }
```

修改 WebMon.cfg 

```
[default]
log_file=web-monitor.log
# Panic 0
# Fatal 1
# Error 2
# Warn 3
# Info 4
# Debug 5
log_level=5

pushurl=http://127.0.0.1:1988/v1/push
#推 agent 到:1988/v1/push
#推 transfer 到:6060/api/push

interval=60

#自定义endpoint
#留空则默认为主机名
endpoint=

[nginx]
#enable 1,disable 0
enabled=1

staturl=http://127.0.0.1/status
```

