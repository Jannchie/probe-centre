package centre

import (
	"github.com/gin-gonic/gin"
	"github.com/jannchie/probe/centre/collector"
	"github.com/jannchie/probe/centre/middleware"
	"github.com/jannchie/probe/centre/tasker"
	. "github.com/jannchie/probe/common"
)

func InitCentre() *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.Cors)
	if err := InitDB(); err != nil {
		panic(err.Error())
	}
	tasker.Init(engine)
	collector.Init(engine)
	return engine
}
