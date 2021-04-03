package router

import (
	"net/http"

	service "github.com/Jannchie/probe-centre/controller"

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
	r.GET("/ping", service.Ping).
		POST("/user", service.CreateUser).
		GET("/user", service.GetUser).
		GET("/token", service.Login).
		POST("/session", service.Login)

	r.Group("/user").
		Use(AuthRequired).
		PUT("/", service.UpdateUser).
		GET("/me", service.GetMe).
		PUT("/token", service.RefreshToken)

	r.Group("task").
		Use(AuthRequired).
		GET("/", service.GetTask).
		POST("/", service.PostTask).
		GET("/stats", service.ListTaskStats)

	r.
		Use(AuthRequired).
		POST("/data", service.PostRaw).
		GET("/ws", service.WsHandler)
}
