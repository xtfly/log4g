log4g is a logging framework that's interface like java log4j. it's designed for configurable and pluggable logging.

--------

## Brief Introduce

core concepts:

 - Formatter：the logger formatter which is identified by `name` attribute, and `type` attribute represents the type of formatter.
 - Output: the logger output target which is identified by `name` attribute, and `type` attribute represents the type of the output target. An output is associated with a formatter.
 - Logger：the logger instances which is identified by `name` attribute. Each logger can be associated with multiple outputs.

## logger

the logger name separating by `/` and `.`, for example, the logger named `module` is the parent logger named `module/submodule`. When the child logger is not set, it' configuration of Output and Level will inherit the parent logger configuration.

Logger level are: All<Trace<Debug<Info<Warn<Error<Critical<Off

The output level of one logger can be configured in the configuration file without case discrimination.

## formatter

the layout of default Formatter is to parse `%{verb}` format string, the type attribute of it is `text`.

the verb format：`%{verbName:fmtstr}`:
 - `verbName` is a field name 
 - `fmtstr` is a format output definitions for the field, it supports standard Go format lattices

for example: `%{pid:05d}`, `05d` means that at least five characters are output, supplemented by `0` . 

anothe example: `%{module}|%{lvl:5s}>>%{msg}`
 
all support verb as follow: 

 - %{pid}: The process id (int)
 - %{program}: The basename of os.Args[0] (string)
 - %{module}: The logger name
 - %{msg}: The content by using Debug(...) or Debugf(...) methods or others of a logger
 - %{level}: The uppercase loglevel name, eg. DEBUG
 - %{lvl}: The uppercase short log level name, eg. DBG
 - %{line}: The line number
 - %{longfile}: The full file path，eg. /a/b/c/d.go
 - %{shortfile}: The file basename, eg. d.go
 - %{longpkg}: The full package name, eg. github.com/xtfly/log4g
 - %{shortpkg}: The package basename, eg. log4g
 - %{longfunc}: The full function name, eg. littleEndian.PutUint32
 - %{shortfunc}: The base function name, eg. PutUint32
 - %{time}: The time when log occurred，eg. %{time:2006-01-02T15:04:05.999Z-07:00}
 - %{xxx}: When using the WithCtx or WithFields method of a logger, `xxx` represents searching for content from the list of output fields.

## output

 **TBD**

## config

The configuration file format supports YAML and JSON formats. For definitions, refer to the `Config` structure in [logger_api.go](api/logger_api.go).

Format:

```
formats:
  - name: f1     # Name of format for output reference
    type: text   # Currently only text
    layout: "%{time} %{level} %{module} %{pid:6d} >> %{msg} (%{longfile}:%{line}) \n"
```

Output:

```
outputs:
  - name: c1          # Name of output for logger reference
    type: console     # Ouput log content into console
    format: f1        # Referenced formatter name
    #async: true      # Whether to start asynchrony ouput log content
    #queue_size: 100  # The length of the queue when enable asynchronous
    #batch_num: 10    # Batch 10 items submitted to the target together when enable asynchronous
    #threshold: info
  - name: r1
    type: size_rolling_file # The type of rolling 
    format: f1
    file: log/rf.log   # The current written output file name
    file_perm: 0640    # The file permissions being written
    back_perm: 0550    # The file permissions that have been backup rolling
    dir_perm: 0750     # The direction permissions
    size: 1M           # When this value is exceeded, make a backup rolling
    backups: 5         # The number of backup rolling
    #async: true       # Whether to start asynchrony ouput log content
    #queue_size: 100   # The length of the queue when enable asynchronous
    #batch_num: 10     # Batch 10 items submitted to the target together when enable asynchronous
    #threshold: info
  - name: r2
    type: time_rolling_file # The type of rolling
    format: f1
    file: log/rf2.log
    file_perm: 0640
    back_perm: 0550
    dir_perm: 0750
    pattern: 2006-01-02   # The date format for backup rolling
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

	_ := log.GetManager().LoadConfigFile("log4g.yaml")

	dlog := log.GetLogger("name")
	dlog.Debug("message")

}
```

## Develop

 - extend Formatter
 - extend Output

## Task List

- [x] Logger formwork
- [x] Formatter
  - [x] Text: parse %{verb} layout
- [x] Output
  - [x] Console
     - [x] sync
     - [x] async
  - [x] Rolling file by size
     - [x] sync
     - [x] async
     - [ ] backup and compress [not test]
  - [x] Rolling file by date
     - [x] sync
     - [x] async
     - [ ] backup and compress [not test]
  - [x] Syslog
     - [x] sync
- [x] validate configuration parameters
- [ ] more test case

## Thanks

Thanks to the opensource project for this project inherits some code of them:

 - [seelog](https://github.com/cihub/seelog)
 - [go-logging](https://github.com/op/go-logging)