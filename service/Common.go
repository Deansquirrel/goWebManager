package service

import (
	"context"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goWebManager/global"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func Start() {
	if global.Args.WebPath != "" && global.Args.WebPort != 0 {
		startWebService(global.Args.WebPath, global.Args.WebPort, global.Args.WebLogLevel)
		return
	}
	cList := strings.Split(global.SysConfig.Iris.PathAndPort, "|")
	cList = goToolCommon.ClearBlock(cList)
	if len(cList) == 0 {
		log.Warn("web info is empty")
		global.Cancel()
		return
	}
	if len(cList)%2 != 0 {
		log.Error("配置数量不为偶数！")
		global.Cancel()
		return
	}

	count := 0
	path := ""
	port := -1
	var err error
	for curr := 0; curr < len(cList); curr = curr + 2 {
		path = cList[curr]
		port, err = strconv.Atoi(cList[curr+1])
		if err != nil {
			log.Error(fmt.Sprintf("get port error: %s,port str is: %s", err.Error(), cList[curr+1]))
		} else {
			startWebService(path, port, global.SysConfig.Iris.LogLevel)
			count = count + 1
		}
	}
	//等待网站启动
	time.Sleep(time.Second * 3)

	if count == 0 {
		log.Info("no web start success")
		global.Cancel()
		return
	}

	expNum := int(len(cList) / 2)
	if expNum == count {
		log.Info(fmt.Sprintf("web strat %d", count))
	} else {
		log.Warn(fmt.Sprintf("web start %d,but exp %d", count, expNum))
	}
}

func startWebService(path string, port int, logLevel string) {
	log.Debug(fmt.Sprintf("start web service,port: %d,path: %s", port, path))
	defer log.Debug(fmt.Sprintf("start web service,port: %d,path: %s Complete", port, path))
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Logger().SetLevel(logLevel)
	app.StaticWeb("/", path)

	go func() {
		_ = app.Run(
			iris.Addr(":"+strconv.Itoa(port)),
			iris.WithoutInterruptHandler,
			iris.WithoutServerError(iris.ErrServerClosed),
			iris.WithOptimizations,
		)
	}()
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch,
			os.Interrupt,
			syscall.SIGINT,
			os.Kill,
			syscall.SIGKILL,
			syscall.SIGTERM,
		)
		select {
		case <-ch:
			stopWebService(app, port, path)
		case <-global.Ctx.Done():
			stopWebService(app, port, path)
		}
	}()
}

func stopWebService(app *iris.Application, port int, path string) {
	log.Debug(fmt.Sprintf("stop web service,port: %d,path: %s", port, path))
	defer log.Debug(fmt.Sprintf("stop web service,port: %d,path: %s Complete", port, path))
	timeout := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_ = app.Shutdown(ctx)
}
