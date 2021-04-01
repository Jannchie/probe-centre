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
	r.GET("/ping", api.Ping)
	r.POST("/user", api.CreateUser)
	r.GET("/user", api.GetUser)
	r.GET("/users", api.ListUserHandle)
	r.GET("/token", api.Login)
	r.POST("/session", api.Login)
	user := r.Group("/user")
	user.Use(AuthRequired)
	user.GET("/probe", api.GetMyProbeList)
	user.PUT("/", api.UpdateUser)
	user.GET("/me", api.GetMe)
	user.PUT("/token", api.RefreshToken)
	r.Use(AuthRequired).
		POST("/stat", api.PostStat).
		POST("/raw", api.PostRaw).
		GET("/task", api.GetTaskHandle).
		POST("/task", api.PostTask).
		GET("/ws", api.WsHandler)
}
