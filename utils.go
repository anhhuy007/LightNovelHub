package main

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
)

func unhex(s string) [16]byte {
	var b [16]byte
	l := len(s)
	for i := 0; i < l; i += 2 {
		out, _ := strconv.ParseInt(s[i:i+2], 16, 16)
		b[i/2] = byte(out)
	}
	return b
}

func getUUID() [16]byte {
	return unhex(strings.ReplaceAll(uuid.New().String(), "-", ""))
}

func passwordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func passwordVerify(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
