package global

import (
	"context"
	"github.com/Deansquirrel/goWebManager/object"
)

const (
	//PreVersion = "0.0.0 Build20190101"
	//TestVersion = "0.0.0 Build20190101"
	Version = "0.0.0 Build20190101"
)

const (
	//http连接超时时间
	HttpConnectTimeout = 30
	//心跳时间间隔
	HeartBeatDuration = 60
)

var Ctx context.Context
var Cancel func()

//程序启动参数
var Args *object.ProgramArgs

//系统参数
var SysConfig *object.SystemConfig

var HasConfig bool
