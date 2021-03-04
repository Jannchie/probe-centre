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

// UserFormID is the form of user by ID.
type UserFormID struct {
	ID uint `form:"ID" binding:"required"`
}

// UserFormToken is the form of user by token.
type UserFormToken struct {
	Token string `form:"Token" binding:"required"`
}

// ListUser is the callback function that returns list of users.
func ListUser(c *gin.Context) {
	var u []model.User
	if c.ShouldBindQuery(&u) == nil {
		res := db.DB.Find(&model.User{}).Limit(10)
		c.JSON(200, res)
	} else {
		c.JSON(400, gin.H{
			"code": -1,
		})
	}
}

// RefreshToken is the callback function that refresh user token.
func RefreshToken(c *gin.Context) {

}

// GetUserByID is the callback function that returns the user.
func GetUserByID(c *gin.Context) {
	var u UserFormID
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

// GetUserByToken is the callback function that returns the user by token.
func GetUserByToken(c *gin.Context) {
	var u UserFormToken
	if err := c.ShouldBindQuery(&u); err != nil {
		c.JSON(400, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	user := model.User{}
	if res := db.DB.First(&user, "token = ?", u.Token); res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  res.Error.Error(),
		})
		return
	}
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

// GetMe is the callback function to get user's information.
func GetMe(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, user)
}
