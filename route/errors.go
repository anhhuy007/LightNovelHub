package route

type ErrorCode uint32

type ErrorJSON struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

const (
	BadInput ErrorCode = iota
	UserNotFound
	WrongPassword
	BadPassword
	BadUsername
	BadDeviceName
	UserAlreadyExists
	InvalidLanguageFormat
	TitleTooLong
	TaglineTooLong
	DescriptionTooLong
)

var message = [...]string{
	"Bad input",
	"User not found",
	"Wrong password",
	"Bad password, password must contains more than 4 letters and less than 72 letters",
	"Bad username, username must contains less than 255 letters",
	"Bad device name, device name must contains less than 32 letters",
	"User already exists, consider login or use another username",
	"Invalid Language, use ISO 639-3",
	"Title too long, title must contains less than 255 letters",
	"Tag line too long, tag line must contains less than 255 letters",
	"Description too long, description must contains less than 5000 letters",
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
