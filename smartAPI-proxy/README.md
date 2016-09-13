## smartAPI-proxy

smartAPI-proxy 负责转发内网的 smartAPI 接口请求


#### 使用方式

配置文件请参照cfg.example.json，修改该文件名为cfg.json

```
{
	"debug": true,
	"smartAPI": {
		"enabled": true,
		"url": "http://smarteye.51idc.com"
	},
    "http": {
        "enabled": true,
        "listen": ":5678"
    }
}


```

#### smartAPI 代理接口
```
curl -d '{"endpoint":"test","version":"0.0.1"}' http://127.0.0.1:5678/api/service/version

```

#### http 信息维护接口

```
curl http://127.0.0.1:5678/health
正常则返回 ok

curl http://127.0.0.1:5678/version
返回版本

curl http://127.0.0.1:5678/workdir
返回工作目录
 
curl http://127.0.0.1:5678/config
返回配置
```

#### 源码安装

```
cd $GOPATH/src/github.com/51idc/service-monitor/smartAPI-proxy/
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
