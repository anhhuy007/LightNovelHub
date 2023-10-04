package repo

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"time"
)

type Database struct {
	db              *sqlx.DB
	timeoutDuration time.Duration
}

func NewDatabase(db *sqlx.DB, timeoutDuration time.Duration) Database {
	return Database{
		db:              db,
		timeoutDuration: timeoutDuration,
	}
}

func (db *Database) countUserFollowers(userID []byte) int {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	follows := 0
	err := db.db.GetContext(
		ctx,
		&follows,
		"SELECT COUNT(*) FROM follows_user WHERE to_id = ?",
		userID,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return 0
	}
	return follows
}

func (db *Database) countUserFollows(fromID []byte) int {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	follows := 0
	err := db.db.GetContext(
		ctx,
		&follows,
		"SELECT COUNT(*) FROM follows_user WHERE from_id = ?",
		fromID,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return 0
	}
	return follows
}

func (db *Database) countComments(toID []byte) int {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	comments := 0
	err := db.db.GetContext(
		ctx,
		&comments,
		"SELECT COUNT(*) FROM comments WHERE to_id = ?",
		toID,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return 0
	}
	return comments
}
