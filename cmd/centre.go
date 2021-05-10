package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jannchie/probe/centre/collector"
	"github.com/jannchie/probe/centre/middleware"
	"github.com/jannchie/probe/centre/tasker"
	common2 "github.com/jannchie/probe/common"
)

func Init() *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.Cors)
	common2.InitDB()
	tasker.Init(engine)
	collector.Init(engine)
	return engine
}

func main() {
	engine := Init()
	addr := "localhost:12000"
	go collector.RunBaseCollector()
	_ = engine.Run(addr)
}
