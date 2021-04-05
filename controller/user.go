package controller

import (
	"errors"
	"net/http"

	"github.com/Jannchie/probe-centre/service"
	"github.com/Jannchie/probe-centre/util"

	"github.com/Jannchie/probe-centre/constant/msg"

	"github.com/Jannchie/probe-centre/constant/code"
	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
	"github.com/Jannchie/probe-centre/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateUser is the callback function that create a user.
func CreateUser(c *gin.Context) {
	type SignUpForm struct {
		Password string `form:"Password" binding:"required"`
		Mail     string `form:"Mail" binding:"required"`
	}
	var form SignUpForm

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code.FAILED,
			"msg":  err.Error(),
		})
		return
	}

	user := model.User{}
	if count := repository.User.CountByMail(form.Mail); count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Mail already been used.",
		})
		return
	}
	newUUID, _ := uuid.NewUUID()
	user.Token = newUUID.String()

	key, salt := service.GenerateKeyAndSalt(form.Password)
	user.Key = key
	user.Salt = salt
	user.Name = form.Mail
	user.Mail = form.Mail
	repository.User.Create(&user)
	c.JSON(http.StatusOK, user)
}

// UserFormID is the form of user by ID.
type UserFormID struct {
	ID uint `form:"ID" binding:"required"`
}

func GetUser(c *gin.Context) {
	// UserFormToken is the form of user by token.
	var form struct {
		Token string `form:"Token"`
		ID    int    `form:"ID"`
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		util.ReturnError(c, err)
		return
	}
	if form.Token != "" {
		GetUserByToken(c, form.Token)
	} else if form.ID != 0 {
		GetUserByID(c, form.ID)
	} else {
		util.ReturnError(c, errors.New(""))
		return
	}
}

// GetUserByID is the callback function that returns the user.
func GetUserByID(c *gin.Context, id int) {
	user := model.User{}
	if res := db.DB.Take(&user, id); res.Error != nil {
		if util.ShouldReturn(c, res.Error) {
			return
		}
	}
	c.JSON(http.StatusOK, user)
}

// GetUserByToken is the callback function that returns the user by token.
func GetUserByToken(c *gin.Context, token string) {
	user := model.User{}
	if res := db.DB.First(&user, "token = ?", token); res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  res.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUserForm is the form for user info update.
type UpdateUserForm struct {
	Name string `form:"Name"`
}

// UpdateUser is the callback function that update the user.
func UpdateUser(c *gin.Context) {
	var u UpdateUserForm
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	user := model.User{}
	db.DB.Model(&user).Where("token", c.Request.Header.Get("token")).Updates(model.User{Name: u.Name})
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": msg.OK, "data": user})
}

// GetMe is the callback function to get user's information.
func GetMe(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, user)
}

// RefreshToken is the callback function that refresh user token.
func RefreshToken(c *gin.Context) {
	token := c.Request.Header.Get("token")
	var user model.User
	newUUID, _ := uuid.NewUUID()
	res := db.DB.Take(&user, "token = ?", token).Update("token", newUUID.String())
	if util.ShouldReturn(c, res.Error) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "data": user.Token, "msg": msg.OK})
}

// Login ist the callback function for login
func Login(c *gin.Context) {

	var form model.LoginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	user, err := service.LoginByForm(form)
	if util.ShouldReturn(c, err) {
		return
	}
	c.JSON(http.StatusOK, user)
}
