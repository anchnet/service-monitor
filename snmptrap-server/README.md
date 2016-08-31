## snmptrap-server

snmptrap-server 通过监听 snmptrap 报文，能够捕获交换机 端口 up/down 的 snmptrap 信息，并直接发送告警。以起到实时响应的效果。
暂时只支持 v2c

#### 开启 snmptrap 
仅供参考
#### 思科
全局配置
```
snmp-server community public RO    #查询的 snmp 配置
snmp-server trap-source Vlan1      #snmp trap 上报时所使用的源地址
snmp-server enable traps snmp linkdown linkup  #开启 linkdown 和 linkup 的 snmptrap
snmp-server host x.x.x.x version 2c public   #将 snmptrap 信息发送到 x.x.x.x,community 是 public，版本 v2c
```
端口配置
当全局开启 snmptrap 时，默认端口也开启 snmptrap。如果不希望所有的端口都发送 trap 信息，则可以关闭不需要监控的端口
进入端口
```
no snmp trap link-status
```

#### 华为/华三
全局配置
```
 snmp-agent
 snmp-agent local-engineid 000007DB7F0000010000562D
 snmp-agent community read cipher %$%$w7ss4!16P;'XF!3Hrsd2SPG>%$%$      #查询的 snmp 配置
 snmp-agent sys-info version v2c
 snmp-agent target-host trap address udp-domai x.x.x.x params securityname public v2c  #将 snmptrap 信息发送到 x.x.x.x,community 是 public,版本 v2c
 snmp-agent trap source Vlanif1 #snmp trap 上报时所使用的源地址
 snmp-agent trap enable feature-name IFNET trap-name linkDown
 snmp-agent trap enable feature-name IFNET trap-name linkUp  #开启 linkdown 和 linkup 的 snmptrap
```
端口配置
当全局开启 snmptrap 时，默认端口也开启 snmptrap。如果不希望所有的端口都发送 trap 信息，则可以关闭不需要监控的端口
进入端口
```
undo enable snmp trap updown
```
#### 使用方式


配置文件请参照cfg.example.json，修改该文件名为cfg.json

```
{
	"debug": true,
	"queue": {
        "sms": "/sms",                  #需和sender的配置一致
        "mail": "/mail"
    },
	"user":{
		"sms":["1381230123","123457890"],
		"mail":["example@exmple.com","1112333238@qq.com"]	
	},
    "redis": {
        "addr": "127.0.0.1:6379",         #需与sender是同一个redis
        "maxIdle": 5
    },
	"trapServer":{
		"trapcommunity":"public",           #snmptrap 上报时携带的 community
		"community":"public",               #snmp查询用的community
		"listen":"0.0.0.0:162"
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
cd $GOPATH/src/github.com/51idc/service-monitor/snmptrap-server/
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

