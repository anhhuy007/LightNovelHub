package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID          []byte         `json:"id"`
	Username    string         `json:"username"`
	Displayname string         `json:"displayName"`
	Password    []byte         `json:"password"`
	Email       sql.NullString `json:"email"`
	Image       string         `json:"image"`
	CreatedAt   time.Time      `json:"created_at"  db:"created_at"`
}

type Novel_status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Novel struct {
	ID          []byte    `json:"id"`
	Title       string    `json:"title"`
	Tagline     string    `json:"tagline"`
	Description string    `json:"description"`
	Author      []byte    `json:"author"`
	Image       string    `json:"image"`
	Language    string    `json:"language"`
	CreateAt    time.Time `json:"createAt"    db:"created_at"`
	UpdateAt    time.Time `json:"updateAt"    db:"updated_at"`
	TotalRating int       `json:"totalRating" db:"total_rating"`
	RateCount   int       `json:"rateCount"   db:"rate_count"`
	Views       int       `json:"views"`
	Clicks      int       `json:"clicks"`
	Adult       bool      `json:"adult"`
	StatusID    int       `json:"statusID"    db:"status_id"`
	Visibility  int       `json:"visibility"`
}

type Tag struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreateAt    time.Time `json:"createAt"    db:"created_at"`
}

type Novel_tags struct {
	NovelID []byte `json:"novelId" db:"novel_id"`
	TagID   int    `json:"tagId"   db:"tag_id"`
}

type Volume struct {
	ID          []byte    `json:"id"`
	NovelID     []byte    `json:"novelId"     db:"novel_id"`
	Title       string    `json:"title"`
	Tagline     string    `json:"tagline"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreateAt    time.Time `json:"createAt"    db:"created_at"`
	UpdateAt    time.Time `json:"updateAt"    db:"updated_at"`
	Views       int       `json:"views"`
	Visibility  int       `json:"visibility"`
}

type Chapter struct {
	ID         []byte    `json:"id"`
	VolumeID   []byte    `json:"volumeId"   db:"volume_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreateAt   time.Time `json:"createAt"   db:"created_at"`
	UpdateAt   time.Time `json:"updateAt"   db:"updated_at"`
	Views      int       `json:"views"`
	Visibility int       `json:"visibility"`
}

type Comment struct {
	ID       []byte    `json:"id"`
	ToID     []byte    `json:"toId"     db:"to_id"`
	UserID   []byte    `json:"userId"   db:"user_id"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"createAt" db:"created_at"`
	UpdateAt time.Time `json:"updateAt" db:"updated_at"`
}

type Image struct {
	ID       int       `json:"id"`
	UserID   []byte    `json:"userId"   db:"user_id"`
	NovelID  []byte    `json:"novelId"  db:"novel_id"`
	Url      string    `json:"url"`
	CreateAt time.Time `json:"createAt" db:"created_at"`
}

type Report struct {
	ID        int       `json:"id"`
	UserID    []byte    `json:"userId"    db:"user_id"`
	ToID      []byte    `json:"toId"      db:"to_id"`
	ReasonID  int       `json:"reasonId"  db:"reason_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type ReportReason struct {
	ID     int    `json:"id"`
	Reason string `json:"reason"`
}

type Session struct {
	ID         []byte    `db:"id"`
	UserID     []byte    `db:"user_id"`
	ExpireAt   time.Time `db:"expires_at"`
	DeviceName string    `db:"device_name"`
}
