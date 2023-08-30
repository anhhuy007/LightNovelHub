package main

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
	"time"
	"unicode/utf8"
)

type Database struct {
	db *sqlx.DB
}

func isUsernameValid(db *Database, username string) bool {
	return !(len(username) == 0 || utf8.RuneCountInString(username) > 63)
}

func (db *Database) GetUser(username string) (User, error) {
	var user User
	if isUsernameValid(db, username) {
		return user, errors.New("Invalid username")
	}
	var tmp context.Context
	ctx, cancel := context.WithTimeout(tmp, time.Minute)
	err := db.db.GetContext(ctx, &user, "SELECT * FROM users WHERE username = ?", username)
	cancel()
	return user, err
}

func (db *Database) CreateUser(username, password string) ([16]byte, error) {
	_, err := db.GetUser(username)
	var userId [16]byte
	if err != nil {
		return userId, errors.New("User already exist")
	}
	hashed, err := passwordHash(password)
	if err != nil {
		return userId, err
	}
	userId = getUUID()
	var tmp context.Context
	ctx, cancel := context.WithTimeout(tmp, time.Minute)
	db.db.ExecContext(ctx, "INSERT INTO users (id, username, password) VALUES (?,?,?)", userId, username, hashed)
	cancel()
	return userId, nil
}
