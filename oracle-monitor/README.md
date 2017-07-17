## oracle-Monitor-agent

监控oracle，采集相关信息上报给 open-falcon

version 信息上报给 smartAPI

以 agent 形式运行，提供简单的 http 信息维护接口

#### 上报字段
视版本和配置不同，采集到的 Metric 可能有所差别。

--------------------------------
| Counters | Type |Tag| Notes|
|-----|------|------|------|
|Oracle.alive|GAUGE|/|oracle alive, 1/-1|
|Oracle.Uptime|GAUGE|database=database,instance=instance|uptime|
|Oracle.tablespace|GAUGE|database=database,instance=instance,tablespace_name=tablespace_name|tablespace usage percent|
|Oracle.sysmetric.User_Calls_Per_Txn|GAUGE|database=database,instance=instance|Calls Per Txn|
|Oracle.sysmetric.Logical_Reads_Per_Sec|GAUGE|database=database,instance=instance|Reads Per Second|
|Oracle.sysmetric.Logical_Reads_Per_Txn|GAUGE|database=database,instance=instance|Reads Per Txn|
|Oracle.sysmetric.Redo_Writes_Per_Sec|GAUGE|database=database,instance=instance|Writes Per Second|
|Oracle.sysmetric.Redo_Writes_Per_Txn|GAUGE|database=database,instance=instance|Writes Per Txn|
|Oracle.sysmetric.Total_Table_Scans_Per_Sec|GAUGE|database=database,instance=instance|Scans Per Second|
|Oracle.sysmetric.Total_Table_Scans_Per_Txn|GAUGE|database=database,instance=instance|Scans Per Txn|
|Oracle.sysmetric.Full_Index_Scans_Per_Sec|GAUGE|database=database,instance=instance|Scans Per Second|
|Oracle.sysmetric.Full_Index_Scans_Per_Txn|GAUGE|database=database,instance=instance|Scans Per Txn|
|Oracle.sysmetric.Execute_Without_Parse_Ratio|GAUGE|database=database,instance=instance|% (ExecWOParse/TotalExec)|
|Oracle.sysmetric.Soft_Parse_Ratio|GAUGE|database=database,instance=instance|% SoftParses/TotalParses|
|Oracle.sysmetric.Host_CPU_Utilization_Ratio|GAUGE|database=database,instance=instance|% Busy/(Idle+Busy)|
|Oracle.sysmetric.DB_Block_Gets_Per_Sec|GAUGE|database=database,instance=instance|Blocks Per Second|
|Oracle.sysmetric.DB_Block_Gets_Per_Txn|GAUGE|database=database,instance=instance|Blocks Per Txn|
|Oracle.sysmetric.Consistent_Read_Gets_Per_Sec|GAUGE|database=database,instance=instance|Blocks Per Second|
|Oracle.sysmetric.Consistent_Read_Gets_Per_Txn|GAUGE|database=database,instance=instance|Blocks Per Txn|
|Oracle.sysmetric.DB_Block_Changes_Per_Sec|GAUGE|database=database,instance=instance|Blocks Per Second|
|Oracle.sysmetric.DB_Block_Changes_Per_Txn|GAUGE|database=database,instance=instance|Blocks Per Txn|
|Oracle.sysmetric.Consistent_Read_Changes_Per_Sec|GAUGE|database=database,instance=instance|Blocks Per Second|
|Oracle.sysmetric.Consistent_Read_Changes_Per_Txn|GAUGE|database=database,instance=instance|Blocks Per Txn|
|Oracle.sysmetric.Database_CPU_Time_Ratio|GAUGE|database=database,instance=instance|% Cpu/DB_Time|
|Oracle.sysmetric.Library_Cache_Hit_Ratio|GAUGE|database=database,instance=instance|% Hits/Pins|
|Oracle.sysmetric.Shared_Pool_Free_Ratio|GAUGE|database=database,instance=instance|% Free/Total|
|Oracle.sysmetric.Executions_Per_Txn|GAUGE|database=database,instance=instance|Executes Per Txn|
|Oracle.sysmetric.Executions_Per_Sec|GAUGE|database=database,instance=instance|Executes Per Second|
|Oracle.sysmetric.Txns_Per_Logon|GAUGE|database=database,instance=instance|Txns Per Logon|
|Oracle.sysmetric.Database_Time_Per_Sec|GAUGE|database=database,instance=instance|CentiSeconds Per Second|
|Oracle.sysmetric.Average_Active_Sessions|GAUGE|database=database,instance=instance|Active Sessions|
|Oracle.sysmetric.Host_CPU_Usage_Per_Sec|GAUGE|database=database,instance=instance|CentiSeconds Per Second|
|Oracle.sysmetric.Cell_Physical_IO_Interconnect_Bytes|GAUGE|database=database,instance=instance|bytes|
|Oracle.sysmetric.Temp_Space_Used|GAUGE|database=database,instance=instance|bytes|
|Oracle.sysmetric.Total_PGA_Allocated|GAUGE|database=database,instance=instance|bytes|
|Oracle.sysmetric.Total_PGA_Used_by_SQL_Workareas|GAUGE|database=database,instance=instance|bytes|
|Oracle.waitmetric.avg_dbtime_wait_1m|GAUGE|database=database,instance=instance,wait_class=wait_class|Percent of database time spent in the wait|
|Oracle.waitmetric.avg_waiter_1m|GAUGE|database=database,instance=instance,wait_class=wait_class|Average waiter count




#### 使用方式


配置文件请参照cfg.example.json，修改该文件名为cfg.json

```
{
  "debug": true,
  "logfile": "oracle.log",
  "hostname": "",
  "db": {
   	"dsn": "system/test@127.0.0.1:1521/orcl",  //需要有 dba 权限
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
    "enabled": false,
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

#### 编译
下载 oracle-client 和 oracle-client-sdk , 准备好 gcc 环境（windows 建议用 MinGW)，准备好 pkg-config 。参考 https://github.com/rana/ora 配置 pkg-config

##### linux 编译
1. 检查下是否有 gcc 环境和 pkgconfig，没有的话安装一下，以 centos 为例
```
yum install pkgconfig gcc
```
2. 安装 oracle-client 和 oracle-client-sdk，官网可下。假定使用 rpm 安装
```
rpm -ivh oracle-instantclient12.2-basic-12.2.0.1.0-1.x86_64.rpm
rpm -ivh oracle-instantclient12.2-devel-12.2.0.1.0-1.x86_64.rpm
```
此时 client 会安装到 ```/usr/lib/oracle/12.2/client64```, sdk 会安装到 ```/usr/include/oracle/12.2/client64```
3. 新建一个目录，用来放 pkgconfig 的配置,并打到环境变量里，例如 ```PKG_CONFIG_PATH=/usr/local/lib/pkgconfig/```
4. 在此目录下增加 ```oci8.pc``` 配置，参考 https://github.com/rana/ora/tree/master/contrib 中的配置
```
prefix=/usr

version=12.2
build=client64

libdir=${prefix}/lib/oracle/${version}/${build}/lib  //这是 client lib 的位置
includedir=${prefix}/include/oracle/${version}/${build} //这是 sdk 的位置

Name: oci8
Description: Oracle database engine
Version: ${version}
Libs: -L${libdir} -lclntsh
Libs.private:
Cflags: -I${includedir}
~                            
```
5. 
```
echo "/usr/lib/oracle/12.2/client64/lib" > /etc/ld.so.conf.d/oracle-client.conf
ldconfig
```
6. 先 ```go get gopkg.in/rana/ora.v4``` 看是 cgo 是否能编译成功
7. 成功 get 下来以后，再编译 oracle-monitor,进入51idc/service-monitor/oracle-monitor目录
```
go get ./...
./control build
```
##### windows 编译
与 linux 类似，需要准备 gcc 环境和 pkgconfig。建议 gcc 用 MinGW
1. 准备环境，[MinGW](http://www.mingw.org/) 和 [pkg-config](https://sourceforge.net/projects/pkgconfiglite/)。把 mingw/bin(mingw64/bin)放入环境变量，同时记得把 pkg-config 放到 mingw/bin 内。
2. 下载 oracle-client和oracle-client-sdk，并解压缩。(非必须：建议同时下一下 sqlplus，并添加一下 oracle-client 的环境变量，用 sqlplus 测试下 oci 调用连接正常
3. 添加环境变量 PKG_CONFIG_PATH ,并在此路径下，添加 oci8.pc 配置
```
# Package Information for pkg-config
prefix=E:/oracle/instantclient_12_2/sdk/

version=12.2
build=client64

libdir=E:/oracle/instantclient_12_2/ //这是 client 的位置
includedir=${prefix}/include //这是 sdk 的位置

glib_genmarshal=glib-genmarshal
gobject_query=gobject-query
glib_mkenums=glib-mkenums

Name: oci8
Description: Oracle database engine
Version: ${version}
Libs: -L${libdir} -loci
Libs.private:
Cflags: -I${includedir}
```
4. 先 ```go get gopkg.in/rana/ora.v4``` 看是 cgo 是否能编译成功
5. 成功 get 下来以后，再编译 oracle-monitor,进入51idc/service-monitor/oracle-monitor目录，编译之
```
go get ./...
go build
```
#### 运行
运行时，需要 oracle 相关 lib 依赖。需要将编译时所使用 oracle oci 的 lib
##### linux
下载打包好的压缩包（其中包含 oracle 12.2 的 oci lib）
```
tar -zxvf falcon-oracle-monitor-0.0.1.tar.gz //解压
echo "$(pwd)/lib" > /etc/ld.so.conf.d/oracle-client.conf //将 lib 库放入 ld 的配置文件中
ldconfig //配置生效
echo 127.0.0.1 ${HOSTNAME} >> /etc/hosts //oracle 要求主机名存在解析，如果本机的主机名没有解析的话，在 hosts 里添加一下
./control start //运行
```
##### windows
1. 解压缩 oracle-monitor-0.0.1.zip
2. 把 lib 目录添加到环境变量PATH中。如果服务器上有部署 oracle 12 的话，不需要此步骤。如果部署的是 11,则需要确保 lib 目录的环境变量位置在 oracle 11 的环境变量目录之前。例如
```
PAHT=C:\Users\Administrator\Desktop\oracle-monitor\lib;C:\app\Administrator\product\11.2.0\dbhome_1\bin;%SystemRoot%\system32;%SystemRoot%;%SystemRoot%\System32\Wbem;%SYSTEMROOT%\System32\WindowsPowerShell\v1.0\
```
3. 运行 oracle-monitor.exe，观察日志是否正确运行
4. 使用 nssm 注册为服务
