package model

import (
	"fmt"
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

	PageSize = 20
)

type NovelStatusID int

const Unknown string = "Unknown"

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
		return Unknown
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
		return Unknown
	}
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

func (order SortOrder) Validate() bool {
	if order != SortOrderDesc && order != SortOrderAsc {
		return false
	}
	return true
}

type OrderBy string

const (
	OrderByCreatedAt OrderBy = "created_at"
	OrderByUpdateAt  OrderBy = "updated_at"
	OrderByViews     OrderBy = "views"
	OrderByTitle     OrderBy = "title"
)

func (order OrderBy) Validate() bool {
	if order != OrderByTitle &&
		order != OrderByCreatedAt &&
		order != OrderByViews &&
		order != OrderByUpdateAt {
		return false
	}
	return true
}

type FiltersAndSortNovel struct {
	SortOrder  SortOrder `db:"sort_order"`
	OrderBy    OrderBy   `db:"order_by"`
	Adult      bool
	Language   string
	Tag        []int
	TagExclude []int `db:"tag_exclude"`
	Search     string
	Page       uint
	FromDate   time.Time `db:"from_date"`
	ToDate     time.Time `db:"to_date"`
	Status     NovelStatusID
}

var DefaultFiltersAndSort = FiltersAndSortNovel{
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
// The generated query will start with AND, if filter with tag,
// please join with table novel tags like this:
// SELECT * FROM (SELECT
//
//	novels.*,
//	GROUP_CONCAT(tag_id) as tag_groupconcat
//
// FROM novels
//
//	RIGHT JOIN novel_tags
//	           ON novels.id = novel_tags.novel_id
//
// GROUP BY novels.id) AS T
// WHERE {query criteria}
// The query should be like this:
// SELECT * FROM novels WHERE 1=1 ...{the generated query here}...
// If the query has it own criteria, the query should be like this:
// SELECT * FROM novels WHERE {query criteria here} ...{the generated query here}...
func (f *FiltersAndSortNovel) ConstructQuery() (string, []interface{}) {
	res := ""
	if f.Adult == false {
		res += " AND novels.adult IS FALSE"
	}
	if f.Status.String() != Unknown {
		res += " AND novels.status = :status"
	}
	if f.Search != DefaultFiltersAndSort.Search {
		res += " AND (MATCH (novels.title) AGAINST (:search IN NATURAL LANGUAGE MODE))"
	}
	if f.Language != DefaultFiltersAndSort.Language {
		res += " AND novels.language LIKE :language"
	}
	if f.FromDate != DefaultFiltersAndSort.FromDate {
		res += " AND novels.created_at >= :from_date"
	}
	if f.ToDate != DefaultFiltersAndSort.ToDate {
		res += " AND novels.created_at <= :to_date"
	}
	for _, tag := range f.Tag {
		res += fmt.Sprintf(" AND FIND_IN_SET(%v, tag_groupconcat)", tag)
	}
	for _, tag := range f.TagExclude {
		res += fmt.Sprintf(" AND FIND_IN_SET(%v, tag_groupconcat) = 0", tag)
	}
	// Maybe a redundant, but who knows?
	if !f.SortOrder.Validate() {
		f.SortOrder = DefaultFiltersAndSort.SortOrder
	}
	res += fmt.Sprintf(" ORDER BY :order_by %v, novels.id ASC", f.SortOrder)
	res += fmt.Sprintf(" LIMIT %v OFFSET %v", PageSize*f.Page, PageSize*(f.Page-1))
	resQuery, args, err := sqlx.Named(res, f)
	if err != nil {
		log.Error(err)
		return "", nil
	}
	//log.Debug(resQuery)
	return resQuery, args
}
