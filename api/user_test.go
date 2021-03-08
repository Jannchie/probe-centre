package api

import (
	"crypto/rand"
	"testing"
	"time"

	"golang.org/x/crypto/argon2"
)

func TestTimeConsuming(t *testing.T) {
	start := time.Now()
	// bcrypt.GenerateFromPassword([]byte("testpassword"), cost)
	key := argon2.IDKey([]byte("testpassword"), []byte("testpassword"), 3, 32*1024, 4, 32)
	t.Logf("cost: %d, duration: %v\n", 1, time.Since(start))
	t.Logf(string(key))
}

func TestGenerateKeyAndSalt(t *testing.T) {
	len := 128
	salt := make([]byte, len)
	t.Log(salt)
	rand.Read(salt)
	t.Log(salt)
	t.Log(string(salt))
}
