package util

import (
	"errors"

	"github.com/Jannchie/probe-centre/model"
	"github.com/gin-gonic/gin"
)

// GetUserFromCtx is a helper function to get the user form gin context.
func GetUserFromCtx(c *gin.Context) (model.User, error) {
	if u, exists := c.Get("user"); exists {
		if user, ok := u.(model.User); ok {
			return user, nil
		}
	}
	return model.User{}, errors.New("get user failed")
}
