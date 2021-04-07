package router

import (
	"log"
	"net/http"
	"time"

	"github.com/Jannchie/probe-centre/util"

	"github.com/Jannchie/probe-centre/constant/code"
	"github.com/Jannchie/probe-centre/controller"
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

func RecordIP(c *gin.Context) {
	c.Next()
	ip := c.ClientIP()
	path := c.FullPath()
	record := model.IPRecord{
		IP:   ip,
		Path: path,
		Time: time.Now().UTC(),
	}
	db.DB.Create(&record)
	log.Println(ip)
	u, err := util.GetUserFromCtx(c)
	if err == nil {
		db.DB.Model(&u).Update("ip", c.ClientIP())
	}
}

// InitRouter init the router
func InitRouter(r *gin.Engine) {
	r.Use(RecordIP)

	r.GET("/ping", controller.Ping).
		POST("/user", controller.CreateUser).
		GET("/user", controller.GetUser).
		GET("/token", controller.Login).
		POST("/session", controller.Login)

	r.Group("/user").
		Use(AuthRequired).
		PUT("/", controller.UpdateUser).
		GET("/me", controller.GetMe).
		PUT("/token", controller.RefreshToken)

	r.Group("task").
		Use(AuthRequired).
		GET("/", controller.GetTask).
		POST("/", controller.PostTask).
		GET("/stats", controller.ListTaskStats)

	r.
		Use(AuthRequired).
		POST("/data", controller.PostRaw).
		GET("/ws", controller.ClientWebSocketHandle)
}
