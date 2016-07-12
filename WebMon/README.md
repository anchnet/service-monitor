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
|Nginx.ServerReading|/|GAUGE|读的 worker 数|
|Nginx.ServerWaiting|/|GAUGE|等待请求的 worker 数|
|Nginx.ServerWriting|/|GAUGE|写的 worker 数|

#### Apache

--------------------------------
| key |  tag | type | note |
|-----|------|------|------|
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
#### Apache
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

修改 WebMon.cfg 

```
[default]
log_file=./web-monitor.log
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

#上报的频率，单位为秒，需和脚本的运行频率设置一致
interval=60

#自定义endpoint
#留空则默认为主机名
endpoint=

[nginx]
#enable 1,disable 0
enabled=1

staturl=http://127.0.0.1/status

[apache]
#enable 1,disable 0
enabled=1

staturl=https://www.apache.org/server-status?auto
```

编译 WebMon

```
git clone https://github.com/51idc/service-monitor.git
cd service-monitor/WebMon/
go build -o WebMon
```


测试 WebMon，假定放在 /opt/webmon 下

```
/opt/webmon/WebMon -c /opt/webmon/WebMon.cfg
```

看下 log 是否运行正常

丢进定时任务 (crontab) 完事