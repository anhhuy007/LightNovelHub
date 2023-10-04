package model

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	IDBinLength = 16
	IDHexLength = IDBinLength * 2

	UsernameMaxLength    = 32
	UsernameMinLength    = 3
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

	DeviceNameMinLength = 0
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
	VisibilityPrivate VisibilityID = 1
	VisibilityPublic  VisibilityID = 2
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

type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

type OrderBy string

const (
	OrderByCreatedAt OrderBy = "created_at"
	OrderByUpdateAt  OrderBy = "updated_at"
	OrderByViews     OrderBy = "views"
	OrderByTitle     OrderBy = "title"
)

type FiltersAndSort struct {
	SortOrder  SortOrder
	OrderBy    OrderBy
	Adult      bool
	Language   string
	Tag        []int
	TagExclude []int
	Search     string
	Page       uint
	FromDate   time.Time
	ToDate     time.Time
	Status     NovelStatusID
}

var DefaultFiltersAndSort = FiltersAndSort{
	SortOrder:  SortOrderDesc,
	OrderBy:    OrderByCreatedAt,
	Adult:      false,
	Language:   "",
	Tag:        []int{},
	TagExclude: []int{},
	Search:     "",
	Page:       1,
	FromDate:   time.Time{},
	ToDate:     time.Time{},
	Status:     NovelStatusID(0),
}

// Return the after WHERE clause to the end in the query.
// The generated query will start with AND, if filter with tag, please join with table novel tags.
// The query should be like this:
// SELECT * FROM novels WHERE 1=1 ...{the generated query here}...
// If the query has it own criteria, the query should be like this:
// SELECT * FROM novels WHERE {query criteria here} ...{the generated query here}...
func (f *FiltersAndSort) ConstructQuery() (string, []interface{}) {
	resQuery, args, err := sqlx.Named(":value", f)
	if err != nil {
		log.Error(err)
		return "", nil
	}
	return resQuery, args
}
