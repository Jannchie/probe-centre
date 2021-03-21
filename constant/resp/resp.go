package resp

import (
	"github.com/Jannchie/probe-centre/constant/code"
	"github.com/Jannchie/probe-centre/constant/msg"
	"github.com/gin-gonic/gin"
)

var (
	OK = gin.H{
		"code": code.OK,
		"msg":  msg.OK,
	}
)
