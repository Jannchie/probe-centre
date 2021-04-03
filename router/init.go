package router

import (
	"net/http"

	"github.com/Jannchie/probe-centre/api"
	"github.com/Jannchie/probe-centre/constant/code"
	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
	"github.com/gin-gonic/gin"
)

// AuthRequired is the function to validate user token.
func AuthRequired(c *gin.Context) {
	token := c.Request.Header.Get("token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": code.FAILED,
			"msg":  "Login required.",
		})
		return
	}
	var user model.User
	res := db.DB.First(&user, "token = ?", token)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": code.FAILED,
			"msg":  res.Error.Error(),
		})
		return
	}
	c.Set("user", user)
	c.Next()
}

// Init initializes the router
func Init() *gin.Engine {
	r := gin.Default()
	InitRouter(r)
	return r
}

// InitRouter init the router
func InitRouter(r *gin.Engine) {
	r.GET("/ping", api.Ping).
		POST("/user", api.CreateUser).
		GET("/user", api.GetUser).
		GET("/token", api.Login).
		POST("/session", api.Login)

	r.Group("/user").
		Use(AuthRequired).
		PUT("/", api.UpdateUser).
		GET("/me", api.GetMe).
		PUT("/token", api.RefreshToken)

	r.Group("task").
		Use(AuthRequired).
		GET("/", api.GetTask).
		POST("/", api.PostTask).
		GET("/stats", api.ListTaskStats)

	r.
		Use(AuthRequired).
		POST("/data", api.PostRaw).
		GET("/ws", api.WsHandler)
}
