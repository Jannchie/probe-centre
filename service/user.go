package service

import (
	"crypto/rand"

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
