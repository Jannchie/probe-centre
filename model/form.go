package model

type LoginForm struct {
	Mail     string `form:"Mail" binding:"required"`
	Password string `form:"Password" binding:"required"`
}
