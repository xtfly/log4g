log4g is a logging framework that's interface like java log4j. it's designed for configurable and pluggable logging.

--------

## Brief Introduce

log4g是一个使用接口与log4j相似的Go语言实现，专为可配置和可扩展而设计的日志框架，它涉及到三个概念：

 - formatter：日志格式化，通过`name`来标识，`type`表示格式化器的类型。
 - output: 日志输出目标，通过`name`来标识，`type`表示输出目标的类型，一个output关联一个formatter。
 - logger：日志实例，通过`name`来标识，每个logger可以关联到多个output。

## logger

logger name支持按`/`与`.`来分割名称。例如名为`module`是`module/submodule`的父Logger,
当子Logger没有设置时，会共享使用父Logger的Output与Level。

logger level支持All<Trace<Debug<Info<Warn<Error<Critical<Off，八个级别，可在配置文件配置输出的级别，配置时不区分大小写。

 **TBD**

## formatter

默认提供按`%{verb}`解析的layout的Formatter，其类型为`format`。

支持verb的格式为：`%{verbName:fmtstr}``，其中`verbName`为字段名称，`fmtstr`为字段格式化输出定义，支持标准的Go的format格式。
如`%{pid:05d}``,其中`05d`表示至少输出5个字符，以0补充。例如输出日志格式为`%{module}|%{lvl:5s}>>%{msg}`。

 - %{pid}: 输出进程ID
 - %{program}：输出程序名
 - %{module}：输出日志名称
 - %{msg}：输出Logger.Debug(...)与Debugf(...)等类似方法的内容
 - %{level}：输出日志级别，大全写，如DEBUG
 - %{lvl}：输出日志级别的三个字母的缩写，大全写，如
 - %{line}：输出打印日志所在行数
 - %{longfile}：输出打印日志所在文件路径，如 /a/b/c/d.go
 - %{shortfile}：输出打印日志所在文件名，如d.go
 - %{longpkg}：输出打印日志所在package全路径名称，如github.com/xtfly/log4g
 - %{shortpkg}：输出打印日志所在package名称，如log4g
 - %{longfunc}：输出打印日志所在函数或方法全名，如littleEndian.PutUint32
 - %{shortfunc}：输出打印日志所在函数或方法名，如PutUint32
 - %{time}：输出当前时间，如%{time:2006-01-02T15:04:05.999Z-07:00}
 - %{xxx}：当使用logger.WithCtx与logger.WithFields接口，xxx表示从输出的字段列表中搜索到内容。

## output

 **TBD**

## config

配置文件格式支持yaml与json格式，定义请参考[api.go](api.go)中`Config`结构体

Format:

```
formats:
  - name: f1     # format的名称，用于output引用
    type: text # 当前只能为text
    layout: "%{time} %{level} %{module} %{pid:6d} >> %{msg} (%{longfile}:%{line}) \n"
```

Output:

```
outputs:
  - name: c1          # output的名称
    type: console     # 表示输出到console
    format: f1        # 引用的formatter名称
    #async: true      # 是否启动异步
    #queue_size: 100  # 异步时，队列的长度
    #batch_num: 10    # 异步时，批量10条一起提交到文件
    #threshold: info
  - name: r1
    type: size_rolling_file # 绕接日志的类型
    format: f1
    file: log/rf.log   # 正在写的文件
    file_perm: 0640    # 正在写的文件权限
    back_perm: 0550    # 已绕接备份的文件权限
    dir_perm: 0750     # 日志目录权限
    size: 1M           # 当超过此值进行绕接备份
    backups: 5         # 绕接备份的个数
    #async: true       # 是否启动异步
    #queue_size: 100   # 异步时，队列的长度
    #batch_num: 10     # 异步时，批量10条一起提交到文件
    #threshold: info
  - name: r2
    type: time_rolling_file # 绕接日志的类型
    format: f1
    file: log/rf2.log
    file_perm: 0640
    back_perm: 0550
    dir_perm: 0750
    pattern: 2006-01-02   # 日期绕接备份的格式
    backups: 5
    #async: true
    #queue_size: 100
    #batch_num: 10
    #threshold: info
  - name: s1
    type: syslog
    format: f1
    prefix: module
```


## usage

```
import "github.com/xtfly/log4g"

func main() {

	_ := log.GetManager().LoadConfigFile("log4g.yml")

	dlog := log.GetLogger("name")
	dlog.Debug("message")

}
```

## Develop

 - extend Formatter
 - extend Output

## Task List

- [x] 日志框架
- [x] Formatter
  - [x] Text: 按%{verb}解析layout
- [x] Output
  - [x] Console输出
     - [x] 同步
     - [x] 异步
  - [ ] 基于大小绕接文件输出
     - [x] 同步
     - [x] 异步
     - [ ] 备份压缩（未验证）
  - [ ] 基于日期绕接文件输出
     - [x] 同步
     - [x] 异步
     - [ ] 备份压缩（未验证）
  - [x] 输出到syslog
     - [x] 同步
- [ ] 配置参数检查
- [ ] 更多的测试覆盖

## Thanks

本项目部分代码继承了[seelog](https://github.com/cihub/seelog)与[go-logging](https://github.com/op/go-logging)的代码，在此表示感谢。
