package model

const (
	IDBinLength = 16
	IDHexLength = 32

	UserNameMaxLength    = 32
	UserNameMinLength    = 3
	PasswordMaxLength    = 72
	PasswordMinLength    = 8
	DisplayNameMaxLength = 64
	DisplayNameMinLength = 3
	EmailMaxLength       = 255
	ImageURLMaxLength    = 255

	TitleMaxLength        = 255
	TaglineMaxLength      = 255
	DescriptionMaxLength  = 5000
	ContentMaxLength      = 16777215 // 16MB 2^24
	CommentMaxLength      = 5000
	ReportReasonMaxLength = 50

	TagNameMaxLength        = 50
	TagNameMinLength        = 2
	TagDescriptionMaxLength = 300

	DeviceNameMaxLength = 255
)

type NovelStatusID int

const (
	StatusOngoing   NovelStatusID = 1
	StatusCompleted NovelStatusID = 2
	StatusDropped   NovelStatusID = 3
)

func (n NovelStatusID) String() string {
	switch n {
	case StatusOngoing:
		return "Ongoing"
	case StatusCompleted:
		return "Completed"
	case StatusDropped:
		return "Dropped"
	default:
		return "Unknown"
	}
}

type VisibilityID int

const (
	VisibilityPrivate VisibilityID = 2
	VisibilityPublic  VisibilityID = 1
)

func (v VisibilityID) String() string {
	switch v {
	case VisibilityPublic:
		return "PUB"
	case VisibilityPrivate:
		return "PRI"
	default:
		return "Unknown"
	}
}
