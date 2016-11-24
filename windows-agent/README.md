falcon-agent
===

This is a windows monitor agent for 51idc

## 


## Installation

It is a golang classic project

```
# set $GOPATH and $GOROOT
git clone https://github.com/51idc/service-monitor/windows-agent
cd agent
go build
start.bat //启动
stop.bat //关闭
```

## Configuration

- heartbeat: heartbeat server rpc address
- transfer: transfer rpc address
- ignore: the metrics should ignore


## 说明
- 数据采集用的是gopsutil（修改了部分代码），具体请查看看https://github.com/shirou/gopsutil
- agent不提供插件、http和缓存功能
- 编译环境需要go-1.5+

