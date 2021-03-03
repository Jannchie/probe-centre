package api

import (
	"net/http"

	"github.com/Jannchie/pyobe-carrier/db"
	"github.com/Jannchie/pyobe-carrier/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateUser is the callback function that create a user.
func CreateUser(c *gin.Context) {
	user := model.User{}
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	var count int64
	if db.DB.Model(&user).Where("mail = ?", user.Mail).Count(&count); count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Mail already been used.",
		})
		return
	}
	uuid, _ := uuid.NewUUID()
	user.Token = uuid.String()
	db.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
	})
}

// UserForm is the form of user
type UserForm struct {
	ID uint `form:"ID" binding:"required"`
}

// ListUser is the callback function that returns list of users.
func ListUser(c *gin.Context) {
	var u UserForm
	if c.ShouldBindQuery(&u) == nil {
		res := db.DB.Find(&model.User{}, u.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, gin.H{
			"code": -1,
		})
	}
}

// GetUser is the callback function that returns the user.
func GetUser(c *gin.Context) {
	var u UserForm
	if err := c.ShouldBindQuery(&u); err != nil {
		c.JSON(400, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	user := model.User{}
	db.DB.First(&user, u.ID)
	c.JSON(200, user)
}

// UpdateUserForm is the form for user info update.
type UpdateUserForm struct {
	ID   uint   `form:"ID" binding:"required"`
	Name string `form:"Name"`
	Mail string `form:"Mail"`
}

// UpdateUser is the callback function that update the user.
func UpdateUser(c *gin.Context) {
	var u UpdateUserForm
	if err := c.ShouldBind(&u); err != nil {
		c.JSON(400, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	user := model.User{}
	db.DB.Model(&user).Where("id", u.ID).Updates(model.User{Name: u.Name, Mail: u.Mail})
	c.JSON(200, gin.H{"code": 1, "msg": "success"})
}
