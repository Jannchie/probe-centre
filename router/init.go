package router

import (
	"net/http"

	"github.com/Jannchie/pyobe-carrier/api"
	"github.com/Jannchie/pyobe-carrier/db"
	"github.com/Jannchie/pyobe-carrier/model"
	"github.com/gin-gonic/gin"
)

// AuthRequired is the function to validate user token.
func AuthRequired(c *gin.Context) {
	token := c.Request.Header.Get("token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": -1,
			"msg":  "Login required.",
		})
		return
	}
	var user model.User
	res := db.DB.First(&user, "token = ?", token)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": -1,
			"msg":  res.Error.Error(),
		})
		return
	}
	c.Set("user", user)
	c.Next()
}

// InitRouter initializes the router
func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", api.Ping)
	r.POST("/user", api.CreateUser)
	r.GET("/user", api.GetUserByToken)
	r.GET("/token", api.Login)
	r.POST("/session", api.Login)
	user := r.Group("/user")
	user.Use(AuthRequired)
	{
		user.GET("probe", api.GetMyProbeList)
		user.PUT("", api.UpdateUser)
		user.GET("/me", api.GetMe)
		user.PUT("/token", api.RefreshToken)
	}

	r.Use(AuthRequired).POST("/stat", api.PostStat)
	return r
}
