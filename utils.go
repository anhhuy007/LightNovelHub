package utils

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
)

func Unhex(s string) [16]byte {
	var b [16]byte
	l := len(s)
	if l > 32 {
		l = 32
	}

	for i := 0; i < l; i += 2 {
		out, _ := strconv.ParseInt(s[i:i+2], 16, 16)
		b[i/2] = byte(out)
	}
	return b
}

func GetUUID() []byte {
	uid := Unhex(strings.ReplaceAll(uuid.New().String(), "-", ""))
	return uid[:]
}

func PasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func PasswordVerify(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
