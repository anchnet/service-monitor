## Redis-Monitor-agent

通过监控 Redis info 命令，采集相关信息上报给 open-falcon

version 信息上报给 smartAPI

以 agent 形式运行，提供简单的 http 信息维护接口

#### 上报字段
视版本和配置不同，采集到的 Metric 可能有所差别。

--------------------------------
| Counters | Type | Notes|
|-----|------|------|
|aof_current_rewrite_time_sec  |GAUGE|当前AOF重写持续的耗时|
|aof_enabled                   |GAUGE| appenonly是否开启,appendonly为yes则为1,no则为0|
|aof_last_bgrewrite_status     |GAUGE|最近一次AOF重写操作是否成功|
|aof_last_rewrite_time_sec     |GAUGE|最近一次AOF重写操作耗时|
|aof_last_write_status         |GAUGE|最近一次AOF写入操作是否成功|
|aof_rewrite_in_progress       |GAUGE|AOF重写是否正在进行|
|aof_rewrite_scheduled         |GAUGE|AOF重写是否被RDB save操作阻塞等待|
|blocked_clients               |GAUGE|正在等待阻塞命令（BLPOP、BRPOP、BRPOPLPUSH）的客户端的数量|
|client_biggest_input_buf      |GAUGE|当前客户端连接中，最大输入缓存|
|client_longest_output_list    |GAUGE|当前客户端连接中，最长的输出列表|
|cluster_enabled               |GAUGE|是否启用Redis集群模式，cluster_enabled|
|connected_clients             |GAUGE|当前已连接的客户端个数|
|connected_clients_pct         |GAUGE|已使用的连接数百分比，connected_clients/maxclients |
|connected_slaves              |GAUGE|已连接的Redis从库个数|
|evicted_keys                  |COUNTER|因内存used_memory达到maxmemory后，每秒被驱逐的key个数|
|expired_keys                  |COUNTER|因键过期后，被惰性和主动删除清理key的每秒个数|
|hz			       |GAUGE|serverCron执行的频率，默认10，表示100ms执行一次，建议不要大于120|
|instantaneous_input_kbps      |GAUGE|瞬间的Redis输入网络流量(kbps)|
|instantaneous_ops_per_sec     |GAUGE|瞬间的Redis操作QPS|
|instantaneous_output_kbps     |GAUGE|瞬间的Redis输出网络流量(kbps)|
|keyspace_hit_ratio            |GAUGE|查找键的命中率|
|keyspace_hits                 |COUNTER|查找键命中的次数|
|keyspace_misses               |COUNTER|查找键未命中的次数|
|latest_fork_usec              |GAUGE|最近一次fork操作的耗时的微秒数(BGREWRITEAOF,BGSAVE,SYNC等都会触发fork),当并发场景fork耗时过长对服务影响较大|
|loading		       |GAUGE|标志位，是否在载入数据文件|
|master_repl_offset            |GAUGE|master复制的偏移量，除了写入aof外，Redis定期为自动增加|
|mem_fragmentation_ratio       |GAUGE|内存碎片率，used_memory_rss/used_memory|
|pubsub_channels               |GAUGE|目前被订阅的频道数量|
|pubsub_patterns               |GAUGE|目前被订阅的模式数量|
|rdb_bgsave_in_progress        |GAUGE|标志位，记录当前是否在创建RDB快照|
|rdb_current_bgsave_time_sec   |GAUGE|当前bgsave执行耗时秒数|
|rdb_last_bgsave_status        |GAUGE|标志位，记录最近一次bgsave操作是否创建成功|
|rdb_last_bgsave_time_sec      |GAUGE|最近一次bgsave操作耗时秒数|
|rdb_last_save_time            |GAUGE|最近一次创建RDB快照文件的Unix时间戳|
|rdb_changes_since_last_save   |GAUGE|从最近一次dump快照后，未被dump的变更次数(和save里变更计数器类似)|
|rejected_connections          |COUNTER|因连接数达到maxclients上限后，被拒绝的连接个数|
|repl_backlog_active           |GAUGE|标志位，master是否开启了repl_backlog,有效地psync(2.8+)|
|repl_backlog_first_byte_offset|GAUGE|repl_backlog中首字节的复制偏移位|
|repl_backlog_histlen          |GAUGE|repl_backlog当前使用的字节数|
|repl_backlog_size             |GAUGE|repl_backlog的长度(repl-backlog-size)，网络环境不稳定的，建议调整大些
|sync_full                     |GAUGE|累计Master full sync的次数;如果值比较大，说明常常出现全量复制，就得分析原因，或调整repl-backlog-size|
|sync_partial_err              |GAUGE|累计Master pysync 出错失败的次数|
|sync_partial_ok               |GAUGE|累计Master psync成功的次数|
|total_commands_processed      |COUNTER|每秒执行的命令数，比较准确的QPS|
|total_connections_received    |COUNTER|每秒新创建的客户端连接数|
|total_net_input_bytes         |COUNTER|Redis每秒网络输入的字节数|
|total_net_output_bytes        |COUNTER|Redis每秒网络输出的字节数|
|uptime	       |GAUGE|Redis运行时长的秒数|
|used_cpu_sys                  |COUNTER|Redis进程消耗的sys cpu|
|used_cpu_user                 |COUNTER|Redis进程消耗的user cpu|
|used_memory                   |GAUGE|由Redis分配的内存的总量，字节数|
|used_memory_lua               |GAUGE|lua引擎使用的内存总量，字节数；有使用lua脚本的注意监控|
|used_memory_pct               |GAUGE|最大内存已使用百分比,used_memory/maxmemory; 存储场景无淘汰key注意监控.(如果maxmemory=0表示未设置限制,pct永远为0)|
|used_memory_peak              |GAUGE|Redis使用内存的峰值，字节数|
|used_memory_rss               |GAUGE|Redis进程从OS角度分配的物理内存，如key被删除后，malloc不一定把内存归还给OS,但可以Redis进程复用|

#### 使用方式


配置文件请参照cfg.example.json，修改该文件名为cfg.json

```
{
	"debug": true,
	"hostname": "",
	"redis":{
		"enabled": true,
		"addr": "127.0.0.1:6379",
		"password": "",
		"db": 0,
		"ignore": {
	        "Redis.loading": true,
			"Redis.aof_enabled": true
	    }
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
cd $GOPATH/src/github.com/51idc/service-monitor/redis-monitor/
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
./falcon-redis-monitor --check
```
