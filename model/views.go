package model

import (
	"time"
)

type UserView struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	Displayname   string    `json:"displayName"`
	Image         string    `json:"image"`
	CreatedAt     time.Time `json:"created_at"    db:"created_at"`
	NovelCount    int       `json:"novelCount"`
	FollowerCount int       `json:"followCount"`
	FollowedCount int       `json:"followedCount"`
}

type UserMetadata struct {
	Username    string `json:"username"`
	Displayname string `json:"displayname"`
	Email       string `json:"email"`
	Image       string `json:"image"`
}

type UserMetadataSmall struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Displayname string `json:"displayName"`
	Image       string `json:"image"`
}

type TagView struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type NovelView struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Tagline     string            `json:"tagline"`
	Description string            `json:"description"`
	Author      UserMetadataSmall `json:"author"`
	Image       string            `json:"image"`
	Language    string            `json:"language"`
	CreateAt    time.Time         `json:"createAt"    db:"created_at"`
	UpdateAt    time.Time         `json:"updateAt"    db:"updated_at"`
	TotalRating int               `json:"totalRating" db:"total_rating"`
	RateCount   int               `json:"rateCount"   db:"rate_count"`
	Views       int               `json:"views"`
	Clicks      int               `json:"clicks"`
	Adult       bool              `json:"adult"`
	Status      string            `json:"status"      db:"status"`
	Visibility  string            `json:"visibility"`
	Tags        []TagView         `json:"tags"`
	Volumes     int               `json:"volumes"`
	FollowCount int               `json:"followCount"`
}

type NovelMetadata struct {
	Title       string        `json:"title"`
	Tagline     string        `json:"tagline"`
	Description string        `json:"description"`
	Image       string        `json:"image"`
	Language    string        `json:"language"`
	Author      []byte        `json:"-"`
	Visibility  VisibilityID  `json:"visibility"`
	Status      NovelStatusID `json:"status"`
}

type NovelMetadataSmall struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Tagline     string            `json:"tagline"`
	Description string            `json:"description"`
	Author      UserMetadataSmall `json:"author"`
	Image       string            `json:"image"`
	Language    string            `json:"language"`
	TotalRating int               `json:"totalRating" db:"total_rating"`
	RateCount   int               `json:"rateCount"   db:"rate_count"`
	Adult       bool              `json:"adult"`
	Status      string            `json:"status"      db:"status_id"`
	Visibility  string            `json:"visibility"`
	Views       int               `json:"views"`
}
