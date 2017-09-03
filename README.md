log4g is a logging framework that's interface like java log4j.

--------

## Introduce

log4g是一个使用方式上参考log4j的Go语言实现，它涉及到三个概念：

 - formatter：日志格式化，通过`id`来标识，`name`表示格式化器的类型。
 - output: 日志输出目标，通过`id`来标识，`name`表示输出目标的类型，一个output关联一个formatter。
 - logger：日志实例，通过`name`来标识，每个logger可以关联到多个output。


logger name支持按`/`与`.`来分割名称。例如名为`module`是`module/submodule`的父Logger,
当子Logger没有设置时，会共享使用父Logger的Output与Level。


## logger

 **TBD**

## formatter

提供按`%{verb}`解析的layout的Formatter，例如输出日志格式为`%{module}|%{lvl}>>%{msg}`。

支持verb的格式为：`%{verbName:fmtstr}``，其中`verbName`为字段名称，`fmtstr`为字段格式化输出定义，支持标准的Go的format格式。
如`%{pid:05d}``,其中`05d`表示至少输出5个字符，以0补充。

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


## usage

```
import "github.com/xtfly/log4g"

func main() {

	_ := log.GetManager().LoadConfig("log4g.yml")

	dlog := log.GetLogger("name")
	dlog.Debug("message")

}
```

## TODO

- [x] 日志框架
- [x] Formatter
  - [x] 默认按%{verb}解析的layout
- [x] Output
  - [x] Console输出
     - [x] 同步
     - [ ] 异步
  - [ ] 基于大小绕接文件输出
     - [x] 同步
     - [ ] 异步
     - [ ] 备份压缩
  - [ ] 基于日期绕接文件输出
     - [x] 同步
     - [ ] 异步
     - [ ] 备份压缩
  - [ ] 输出到syslog
     - [ ] 同步
     - [ ] 异步
- [ ] 配置参数检查

## Thanks

本项目部分代码采用[seelog](https://github.com/cihub/seelog)与[go-logging](https://github.com/op/go-logging)的代码
