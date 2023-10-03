package model

import (
	"encoding/hex"
	"github.com/google/uuid"
	"strings"
)

func GetUUID() []byte {
	uid, _ := hex.DecodeString(strings.ReplaceAll(uuid.New().String(), "-", ""))
	return uid
}
