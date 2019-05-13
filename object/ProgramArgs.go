package object

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strings"
)

import log "github.com/Deansquirrel/goToolLog"

const (
	ArgsFlagInstall   = "install"
	ArgsFlagUninstall = "uninstall"
	ArgsFlagLogStdOut = "stdout"
	ArgsLogLevel      = "logLevel"

	ArgsWebLogLevel = "webloglevel"
	ArgsWebPath     = "path"
	ArgsWebPort     = "port"

	DefaultLogLevel    = "warn"
	DefaultWebLogLevel = "warn"
)

type ProgramArgs struct {
	IsInstall   bool
	IsUninstall bool

	LogStdOut bool
	LogLevel  string

	WebLogLevel string
	WebPath     string
	WebPort     int
}

func (pa *ProgramArgs) Definition() {
	flag.BoolVar(&pa.IsInstall, ArgsFlagInstall, false, "安装服务")
	flag.BoolVar(&pa.IsUninstall, ArgsFlagUninstall, false, "卸载服务")

	flag.BoolVar(&pa.LogStdOut, ArgsFlagLogStdOut, false, "控制台日志输出")
	flag.StringVar(&pa.LogLevel, ArgsLogLevel, "warn", "日志级别（debug|info|warn|error）")

	flag.StringVar(&pa.WebLogLevel, ArgsWebLogLevel, "", "Web日志级别（debug|info|warn|error）")
	flag.StringVar(&pa.WebPath, ArgsWebPath, "", "Web路径")
	flag.IntVar(&pa.WebPort, ArgsWebPort, 0, "Web端口")
}

func (pa *ProgramArgs) Parse() {
	flag.Parse()
}

func (pa *ProgramArgs) Check() error {
	//安装为服务和卸载服务参数不可同时存在
	if pa.IsInstall && pa.IsUninstall {
		return errors.New(fmt.Sprintf("参数 %s 和 %s 不可同时存在", ArgsFlagInstall, ArgsFlagUninstall))
	}
	pa.LogLevel = strings.Trim(pa.LogLevel, " ")
	if pa.LogLevel != "" &&
		pa.LogLevel != "debug" &&
		pa.LogLevel != "info" &&
		pa.LogLevel != "warn" &&
		pa.LogLevel != "error" {
		log.Warn(fmt.Sprintf("arg log level format to default: %s", DefaultLogLevel))
	}

	pa.WebLogLevel = strings.Trim(pa.WebLogLevel, " ")
	if pa.WebLogLevel != "" &&
		pa.WebLogLevel != "debug" &&
		pa.WebLogLevel != "info" &&
		pa.WebLogLevel != "warn" &&
		pa.WebLogLevel != "error" {
		log.Warn(fmt.Sprintf("arg web log level format to default: %s", DefaultWebLogLevel))
	}
	return nil
}

func (pa *ProgramArgs) ToString() string {
	d, err := json.Marshal(pa)
	if err != nil {
		log.Warn(fmt.Sprintf("ProgramArgs转换为字符串时遇到错误：%s", err.Error()))
		return ""
	}
	return string(d)
}
