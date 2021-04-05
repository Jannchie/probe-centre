package service

import (
	"bytes"
	"crypto/rand"
	"errors"

	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"

	"golang.org/x/crypto/argon2"
)

func GenerateKeyAndSalt(password string) ([]byte, []byte) {
	l := 128
	salt := make([]byte, l)
	_, _ = rand.Read(salt)
	key := argon2.IDKey([]byte(password), salt[:], 3, 32*1024, 4, 32)
	return key, salt
}
func GenerateKeyWithSalt(password string, salt []byte) []byte {
	key := argon2.IDKey([]byte(password), salt[:], 3, 32*1024, 4, 32)
	return key
}

func LoginByForm(form model.LoginForm) (model.User, error) {
	var user model.User
	if err := db.DB.Take(&user, "mail = ?", form.Mail).Error; err != nil {
		return user, err
	}
	key := GenerateKeyWithSalt(form.Password, user.Salt)
	if res := bytes.Compare(key, user.Key); res == 0 {
		return user, nil
	} else {
		return user, errors.New("wrong password or mail")
	}
}
