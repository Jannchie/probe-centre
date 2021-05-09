package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jannchie/probe/centre/collector"
	"github.com/jannchie/probe/centre/common"
	"github.com/jannchie/probe/centre/middleware"
	"github.com/jannchie/probe/centre/tasker"
)

func Init() *gin.Engine {
	// 获取命令行参数
	engine := gin.Default()
	engine.Use(middleware.Cors)
	common.InitDB()
	tasker.Init(engine)
	collector.Init(engine)
	return engine
}

func main() {
	engine := Init()
	_ = engine.Run(":12000")
}
