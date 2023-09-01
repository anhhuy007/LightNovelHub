package model

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrInvalidUsername  = Error("Invalid username")
	ErrInvalidPassword  = Error("Invalid password")
	ErrWrongPassword    = Error("Wrong password")
	ErrUserAlreadyExist = Error("User already exist")
	ErrUserNotFound     = Error("User not found")
	ErrSessionExpired   = Error("Session expired")
	ErrNovelNotFound    = Error("Novel not found")
	ErrInternal         = Error("Internal server error")
)
