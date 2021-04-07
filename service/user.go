package service

import (
	"bytes"
	"errors"

	"github.com/Jannchie/probe-centre/util"

	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
)

func LoginByForm(form model.LoginForm) (model.User, error) {
	var user model.User
	if err := db.DB.Take(&user, "mail = ?", form.Mail).Error; err != nil {
		return user, err
	}
	key := util.GenerateKeyWithSalt(form.Password, user.Salt)
	if res := bytes.Compare(key, user.Key); res == 0 {
		return user, nil
	} else {
		return user, errors.New("wrong password or mail")
	}
}
