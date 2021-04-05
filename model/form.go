package model

type LoginForm struct {
	Mail     string `form:"Mail" binding:"required"`
	Password string `form:"Password" binding:"required"`
}

type RawDataForm struct {
	Data   interface{} `form:"Data"`
	TaskID uint64      `form:"TaskID"`
	Number uint64      `form:"Number"`
}
