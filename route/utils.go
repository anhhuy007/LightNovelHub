package route

import (
	"Lightnovel/model"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"regexp"
	"unicode"
	"unicode/utf8"
)

func PasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func Unhex(s string) ([]byte, error) {
	return hex.DecodeString(s[:model.IDHexLength])
}

func PasswordVerify(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}

func IsUsernameValid(username string) bool {
	matched, _ := regexp.MatchString(
		fmt.Sprintf("^[a-zA-Z]{%v,%v}$", model.UsernameMinLength, model.UsernameMaxLength),
		username,
	)
	return matched
}

func IsPasswordValid(password string) bool {
	return !(len(password) < model.PasswordMinLength || len(password) > model.PasswordMaxLength)
}

func checkUserMetadata(data model.UserMetadata) (bool, ErrorCode) {
	if IsUsernameValid(data.Username) == false {
		return false, BadUsername
	}

	if utf8.RuneCountInString(data.Displayname) < model.DisplayNameMinLength ||
		utf8.RuneCountInString(data.Displayname) > model.DisplayNameMaxLength {
		return false, BadDisplayname
	}

	for _, ch := range []rune(data.Displayname) {
		if !unicode.IsPrint(ch) {
			return false, BadDisplayname
		}
	}

	if _, err := mail.ParseAddress(data.Email); err != nil ||
		utf8.RuneCountInString(data.Email) > model.EmailMaxLength {
		return false, BadEmail
	}

	return true, BadInput
}
