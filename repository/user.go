package repository

import (
	"github.com/Jannchie/probe-centre/model"
	"gorm.io/gorm"
)

type UserReposity struct {
	DB *gorm.DB
}

func (r UserReposity) CountByMail(mail string) int64 {
	var count int64
	r.DB.Model(&model.User{}).Where("mail = ?", mail).Count(&count)
	return count
}

func (r UserReposity) Create(user *model.User) *gorm.DB {
	return r.DB.Create(user)
}
