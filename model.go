package main

import "time"

type User struct {
	ID        [16]byte  `json:"id"`
	Username  string    `json:"username"`
	Password  [60]byte  `json:"password"`
	Email     string    `json:"email"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created" db:"created_at"`
}

type Novel_status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Novel struct {
	ID          [16]byte  `json:"id"`
	Title       string    `json:"title"`
	Tagline     string    `json:"tagline"`
	Description string    `json:"description"`
	Author      [16]byte  `json:"author"`
	Image       string    `json:"image"`
	CreateAt    time.Time `json:"createAt" db:"created_at"`
	UpdateAt    time.Time `json:"updateAt" db:"updated_at"`
	TotalRating int       `json:"totalRating" db:"total_rating"`
	RateCount   int       `json:"rateCount" db:"rate_count"`
	Views       int       `json:"views"`
	Clicks      int       `json:"clicks"`
	Adult       bool      `json:"adult"`
	StatusID    int       `json:"statusId" db:"status_id"`
}

type Tags struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreateAt    time.Time `json:"createAt" db:"created_at"`
}

type Novel_tags struct {
	NovelID [16]byte `json:"novelId" db:"novel_id"`
	TagID   int      `json:"tagId" db:"tag_id"`
}

type Volume struct {
	ID          [16]byte  `json:"id"`
	NovelID     [16]byte  `json:"novelId" db:"novel_id"`
	Title       string    `json:"title"`
	Tagline     string    `json:"tagline"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreateAt    time.Time `json:"createAt" db:"created_at"`
	UpdateAt    time.Time `json:"updateAt" db:"updated_at"`
	Views       int       `json:"views"`
	Clicks      int       `json:"clicks"`
}

type Chapter struct {
	ID       [16]byte  `json:"id"`
	VolumeID [16]byte  `json:"volumeId" db:"volume_id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"createAt" db:"created_at"`
	UpdateAt time.Time `json:"updateAt" db:"updated_at"`
	Views    int       `json:"views"`
	Clicks   int       `json:"clicks"`
}

type Comment struct {
	ID       [16]byte  `json:"id"`
	ToID     [16]byte  `json:"toId" db:"to_id"`
	UserID   [16]byte  `json:"userId" db:"user_id"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"createAt" db:"created_at"`
	UpdateAt time.Time `json:"updateAt" db:"updated_at"`
}

type Image struct {
	ID       int       `json:"id"`
	UserID   [16]byte  `json:"userId" db:"user_id"`
	NovelID  [16]byte  `json:"novelId" db:"novel_id"`
	url      string    `json:"url"`
	CreateAt time.Time `json:"createAt" db:"created_at"`
}

type Report struct {
	ID        int       `json:"id"`
	UserID    [16]byte  `json:"userId" db:"user_id"`
	ToID      [16]byte  `json:"toId" db:"to_id"`
	ReasonID  int       `json:"reasonId" db:"reason_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type ReportReason struct {
	ID     int    `json:"id"`
	Reason string `json:"reason"`
}

type Session struct {
	ID         [16]byte  `db:"id"`
	UserID     [16]byte  `db:"user_id"`
	ExpireAt   time.Time `db:"expire_at"`
	DeviceName string    `db:"device_name"`
}
