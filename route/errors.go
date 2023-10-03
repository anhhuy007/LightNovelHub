package route

import (
	"Lightnovel/model"
	"fmt"
)

type ErrorCode uint32

type ErrorJSON struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

const (
	BadInput ErrorCode = iota

	// Account related error
	UserNotFound
	WrongPassword
	BadPassword
	BadUsername
	BadDeviceName
	BadDisplayname
	BadEmail
	UserAlreadyExists

	// Novel related error
	InvalidLanguageFormat
	TitleTooLong
	TaglineTooLong
	DescriptionTooLong
)

var message = [...]string{
	"Bad input",
	"User not found",
	"Wrong password",
	fmt.Sprintf(
		"Bad password, password must contains more than %v letters and less than %v letters",
		model.PasswordMinLength,
		model.PasswordMaxLength,
	),
	fmt.Sprintf(
		"Bad username, username must contains more than %v letters and less than %v letters",
		model.UsernameMinLength,
		model.UsernameMaxLength,
	),
	fmt.Sprintf(
		"Bad device name, device name must contains more than %v letters and less than %v letters",
		model.DeviceNameMinLength,
		model.DeviceNameMaxLength,
	),
	fmt.Sprintf(
		"Bad displayname, displayname must contains more than %v letters and less than %v letters",
		model.DisplayNameMinLength,
		model.DisplayNameMaxLength,
	),
	fmt.Sprintf(
		"Bad email, email must contains less than %v letters and must be a valid email",
		model.EmailMaxLength,
	),
	"User already exists, consider login or use another username",
	"Invalid Language, use ISO 639-3",
	fmt.Sprintf("Title too long, title must contains less than %v letters", model.TitleMaxLength),
	fmt.Sprintf(
		"Tagline too long, tagline must contains less than %v letters",
		model.TaglineMaxLength,
	),
	fmt.Sprintf(
		"Description too long, description must contains less than %v letters",
		model.DescriptionMaxLength,
	),
}

func getMessage(code ErrorCode) string {
	return message[code]
}

func buildErrorJSON(code ErrorCode) ErrorJSON {
	return ErrorJSON{
		code,
		getMessage(code),
	}
}
