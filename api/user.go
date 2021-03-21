package api

import (
	"bytes"
	"crypto/rand"
	"net/http"

	"github.com/Jannchie/probe-centre/constant/code"
	"github.com/Jannchie/probe-centre/constant/msg"
	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
	"github.com/Jannchie/probe-centre/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

func generateKeyAndSalt(password string) ([]byte, []byte) {
	len := 128
	salt := make([]byte, len)
	rand.Read(salt)
	key := argon2.IDKey([]byte(password), salt[:], 3, 32*1024, 4, 32)
	return key, salt
}
func generateKeyWithSalt(password string, salt []byte) []byte {
	key := argon2.IDKey([]byte(password), salt[:], 3, 32*1024, 4, 32)
	return key
}

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
	uuid, _ := uuid.NewUUID()
	user.Token = uuid.String()

	key, salt := generateKeyAndSalt(form.Password)
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

// UserFormToken is the form of user by token.
type UserFormToken struct {
	Token string `form:"Token" binding:"required"`
}

// ListUser is the callback function that returns list of users.
func ListUser(c *gin.Context) {
	var u []model.User
	if c.ShouldBindQuery(&u) == nil {
		res := db.DB.Find(&model.User{}).Limit(10)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(400, gin.H{
			"code": -1,
		})
	}
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
	c.JSON(http.StatusOK, user)
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
	c.JSON(http.StatusOK, user)
}

// UpdateUserForm is the form for user info update.
type UpdateUserForm struct {
	Name string `form:"Name"`
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
	db.DB.Model(&user).Where("token", c.Request.Header.Get("token")).Updates(model.User{Name: u.Name})
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": msg.OK})
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
	uuid, _ := uuid.NewUUID()
	db.DB.First(&user, "token = ?", token).Update("token", uuid.String())
	c.JSON(http.StatusOK, gin.H{"code": 1, "data": user.Token, "msg": msg.OK})
}

// Login ist the callback function for login
func Login(c *gin.Context) {
	type LoginForm struct {
		Mail     string `form:"Mail" binding:"required"`
		Password string `form:"Password" binding:"required"`
	}
	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	var user model.User
	if err := db.DB.First(&user, "mail = ?", form.Mail).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	key := generateKeyWithSalt(form.Password, user.Salt)
	if res := bytes.Compare(key, user.Key); res == 0 {
		c.JSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": -1,
			"msg":  "wrong password or mail",
		})
	}
}
