## Redis-Monitor-agent

监控 Mongodb ,通过 serverStatus 命令，采集相关信息上报给 open-falcon

version 信息上报给 smartAPI

以 agent 形式运行，提供简单的 http 信息维护接口

#### 上报字段
视版本和配置不同，采集到的 Metric 可能有所差别。

| Counters | Type | Notes|
|-----|------|------|
|asserts_msg|                              COUNTER |消息断言数/秒                                                 
|asserts_regular|                          COUNTER |常规断言数/秒                                                 
|asserts_rollovers|                        COUNTER |计数器roll over的次数/秒,计数器每2^30个断言就会清零           
|asserts_user|                             COUNTER |用户断言数/秒                                                 
|asserts_warning|                          COUNTER |警告断言数/秒                                                 
|page_faults|                              COUNTER |页缺失次数/秒 
|connections_available|                    GAUGE   |未使用的可用连接数                                            
|connections_current|                      GAUGE   |当前所有客户端的已连接的连接数                                
|connections_used_percent|                 GAUGE   |已使用连接数百分比                                            
|connections_totalCreated|                 COUNTER |创建的新连接数/秒                                             
|globalLock_currentQueue_total|            GAUGE   |当前队列中等待锁的操作数                                      
|globalLock_currentQueue_readers|          GAUGE   |当前队列中等待读锁的操作数                                    
|globalLock_currentQueue_writers|          GAUGE   |当前队列中等待写锁的操作数                                    
|locks_Global_acquireCount_ISlock|         COUNTER |实例级意向共享锁获取次数                                      
|locks_Global_acquireCount_IXlock|         COUNTER |实例级意向排他锁获取次数                                      
|locks_Global_acquireCount_Slock|          COUNTER |实例级共享锁获取次数                                          
|locks_Global_acquireCount_Xlock|          COUNTER |实例级排他锁获取次数                                          
|locks_Global_acquireWaitCount_ISlock|     COUNTER |实例级意向共享锁等待次数                                      
|locks_Global_acquireWaitCount_IXlock|     COUNTER |实例级意向排他锁等待次数                                      
|locks_Global_timeAcquiringMicros_ISlock|  COUNTER |实例级共享锁获取耗时 单位:微秒                                
|locks_Global_timeAcquiringMicros_IXlock|  COUNTER |实例级共排他获取耗时 单位:微秒                                
|locks_Database_acquireCount_ISlock|       COUNTER |数据库级意向共享锁获取次数                                    
|locks_Database_acquireCount_IXlock|       COUNTER |数据库级意向排他锁获取次数                                    
|locks_Database_acquireCount_Slock|        COUNTER |数据库级共享锁获取次数                                        
|locks_Database_acquireCount_Xlock|        COUNTER |数据库级排他锁获取次数                                        
|locks_Collection_acquireCount_ISlock|     COUNTER |集合级意向共享锁获取次数                                      
|locks_Collection_acquireCount_IXlock|     COUNTER |集合级意向排他锁获取次数                                      
|locks_Collection_acquireCount_Xlock|      COUNTER |集合级排他锁获取次数                                          
|opcounters_command|                       COUNTER |数据库执行的所有命令/秒                                       
|opcounters_insert|                        COUNTER |数据库执行的插入操作次数/秒                                   
|opcounters_delete|                        COUNTER |数据库执行的删除操作次数/秒                                   
|opcounters_update|                        COUNTER |数据库执行的更新操作次数/秒                                   
|opcounters_query|                         COUNTER |数据库执行的查询操作次数/秒                                   
|opcounters_getmore|                       COUNTER |数据库执行的getmore操作次数/秒                                
|opcountersRepl_command|                   COUNTER |数据库复制执行的所有命令次数/秒                               
|opcountersRepl_insert|                    COUNTER |数据库复制执行的插入命令次数/秒                               
|opcountersRepl_delete|                    COUNTER |数据库复制执行的删除命令次数/秒                               
|opcountersRepl_update|                    COUNTER |数据库复制执行的更新命令次数/秒                               
|opcountersRepl_query|                     COUNTER |数据库复制执行的查询命令次数/秒                               
|opcountersRepl_getmore|                   COUNTER |数据库复制执行的gtemore命令次数/秒                            
|network_bytesIn|                          COUNTER |数据库接受的网络传输字节数/秒                                 
|network_bytesOut|                         COUNTER |数据库发送的网络传输字节数/秒                                 
|network_numRequests|                      COUNTER |数据库接收到的请求的总次数/秒                                 
|mem_virtual|                              GAUGE   |数据库进程使用的虚拟内存                                      
|mem_resident|                             GAUGE   |数据库进程使用的物理内存                                      
|mem_mapped|                               GAUGE   |mapped的内存,只用于MMAPv1 存储引擎                            
|mem_bits|                                 GAUGE   |64 or 32bit                                                   
|mem_mappedWithJournal|                    GAUGE   |journal日志消耗的映射内存，只用于MMAPv1 存储引擎              
|backgroundFlushing_flushes|               COUNTER |数据库刷新写操作到磁盘的次数/秒                               
|backgroundFlushing_average_ms|            GAUGE   |数据库刷新写操作到磁盘的平均耗时，单位ms                      
|backgroundFlushing_last_ms|               COUNTER |当前最近一次数据库刷新写操作到磁盘的耗时，单位ms              
|backgroundFlushing_total_ms|              GAUGE   |数据库刷新写操作到磁盘的总耗时/秒，单位ms                     
|cursor_open_total|                        GAUGE   |当前数据库为客户端维护的游标总数                              
|cursor_timedOut|                          COUNTER |数据库timout的游标个数/秒                                     
|cursor_open_noTimeout|                    GAUGE   |设置DBQuery.Option.noTimeout的游标数                          
|cursor_open_pinned|                       GAUGE   |打开的pinned的游标数                                          
|wt_cache_used_total_bytes|                GAUGE   |wiredTiger cache的字节数                                      
|wt_cache_dirty_bytes|                     GAUGE   |wiredTiger cache中"dirty"数据的字节数                         
|wt_cache_readinto_bytes|                  COUNTER |数据库写入wiredTiger cache的字节数/秒                         
|wt_cache_writtenfrom_bytes|               COUNTER |数据库从wiredTiger cache写入到磁盘的字节数/秒                 
|wt_concurrentTransactions_write|          GAUGE   |write tickets available to the WiredTiger storage engine      
|wt_concurrentTransactions_read|           GAUGE   |read tickets available to the WiredTiger storage engine       
|wt_bm_bytes_read|                         COUNTER |block-manager read字节数/秒                                   
|wt_bm_bytes_written|                      COUNTER |block-manager write字节数/秒                                  
|wt_bm_blocks_read|                        COUNTER |block-manager read块数/秒                                     
|wt_bm_blocks_written|                     COUNTER |block-manager write块数/秒    

#### 使用方式


配置文件请参照cfg.example.json，修改该文件名为cfg.json

```
{
	"debug": true,
	"hostname": "",
	"mongo":{
		"enabled": true,
		"addr": "127.0.0.1:27017",
		"username": "",
		"password": "",
		"authdb": "",
		"ignore": {
	        "Mongo.backgroundFlushing_flushes": true
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
cd $GOPATH/src/github.com/51idc/service-monitor/mongo-monitor/
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
./falcon-mongo-monitor --check
```
