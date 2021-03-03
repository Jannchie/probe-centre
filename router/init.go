package router

import (
	"github.com/Jannchie/pyobe-carrier/api"
	"github.com/gin-gonic/gin"
)

// InitRouter initializes the router
func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", api.Ping)
	r.POST("/stat", api.PostStat)
	r.POST("/user", api.CreateUser)
	r.PUT("/user", api.UpdateUser)
	r.GET("/user", api.GetUser)
	return r
}
