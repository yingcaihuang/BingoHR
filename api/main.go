package main

import (
	"context"
	"errors"
	"fmt"
	"hr-api/pkg/bus"
	"hr-api/pkg/cache"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"hr-api/models"
	"hr-api/pkg/logging"
	"hr-api/pkg/setting"
	"hr-api/pkg/util"
	"hr-api/routers"
)

func init() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	cache.Init()
	util.Setup()
}

// @title Golang Gin API
// @version 1.0
// @description An example of gin
// @termsOfService https://hr-api
// @license.name MIT
// @license.url https://hr-api/blob/master/LICENSE
func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	// 从service bus消费需要分析的简历
	go func() {
		log.Println("resume worker started")
		if err := bus.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.Println("resume worker error:", err)
			stop() // 真异常才停止程序
		}
	}()

	server.ListenAndServe()

	// If you want Graceful Restart, you need a Unix system and download github.com/fvbock/endless
	//endless.DefaultReadTimeOut = readTimeout
	//endless.DefaultWriteTimeOut = writeTimeout
	//endless.DefaultMaxHeaderBytes = maxHeaderBytes
	//server := endless.NewServer(endPoint, routersInit)
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}
	//
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Printf("Server err: %v", err)
	//}
}
