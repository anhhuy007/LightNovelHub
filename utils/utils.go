package utils

import (
	"encoding/hex"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func Unhex(s string) []byte {
	b, _ := hex.DecodeString(s[:32])
	return b
}

func GetUUID() []byte {
	uid, _ := hex.DecodeString(strings.ReplaceAll(uuid.New().String(), "-", ""))
	return uid
}

func PasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
