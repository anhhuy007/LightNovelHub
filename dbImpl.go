package db

import (
	"Lightnovel"
	"Lightnovel/utils"
	"database/sql"
	"encoding/hex"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
	"time"
	"unicode/utf8"
)

type Database struct {
	db *sqlx.DB
}

func isUsernameValid(username string) bool {
	return !(len(username) == 0 || utf8.RuneCountInString(username) > 63)
}

func isPasswordValid(password string) bool {
	return !(len(password) < 10 || len(password) > 72)
}

func (db *Database) GetUser(username string) (User, error) {
	var user User
	if !isUsernameValid(username) {
		return user, main.ErrInvalidUsername
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	err := db.db.GetContext(
		ctx,
		&user,
		"SELECT id, username, password, email, image, created_at FROM users WHERE username = ?",
		username,
	)
	cancel()
	if err == sql.ErrNoRows {
		return User{}, main.ErrUserNotFound
	} else if err != nil && err != context.Canceled {
		log.Error(err)
		return User{}, main.ErrInternal
	}
	return user, nil
}

func (db *Database) CreateUser(username, password string) ([]byte, error) {
	var userId []byte
	_, err := db.GetUser(username)
	if err == nil {
		return userId, main.ErrUserAlreadyExist
	}

	if !isPasswordValid(password) {
		return userId, main.ErrInvalidPassword
	}
	hashed, err := utils.PasswordHash(password)
	if err != nil {
		log.Debug(err)
		return userId, main.ErrInvalidPassword
	}

	userId = utils.GetUUID()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_, err = db.db.ExecContext(
		ctx,
		"INSERT INTO users (id, username, password) VALUES (?,?,?)",
		userId,
		username,
		hashed,
	)
	cancel()

	if err != nil && err != context.Canceled {
		log.Error(err)
		return userId, main.ErrInternal
	}
	return userId, nil
}

func (db *db.Database) CreateSession(
	userID []byte,
	deviceName string,
) (SessionInfo, error) {
	sessionID := utils.GetUUID()
	expires := time.Now().Add(sessionDuration)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_, err := db.db.ExecContext(
		ctx,
		"INSERT INTO sessions (id, user_id, expires_at, device_name) VALUES (?, ?, ?, ?)",
		sessionID,
		userID,
		expires,
		deviceName,
	)
	cancel()
	if err != nil && err != context.Canceled {
		log.Error(err)
		return SessionInfo{}, ErrInternal
	}
	return SessionInfo{hex.EncodeToString(sessionID[:]), expires}, nil
}

func (db *db.Database) GetSession(sessionID string) (db.Session, error) {
	var session db.Session
	sessionIDByte := utils.Unhex(sessionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	err := db.db.GetContext(
		ctx,
		&session,
		"SELECT id, user_id, expires_at, device_name FROM sessions WHERE id = ?",
		sessionIDByte[:],
	)
	cancel()

	if err == sql.ErrNoRows {
		return db.Session{}, ErrSessionExpired
	} else if err != nil && err != context.Canceled {
		log.Error(err)
		return db.Session{}, ErrInternal
	}
	if time.Now().After(session.ExpireAt) {
		return db.Session{}, ErrSessionExpired
	}

	if session.ExpireAt.Sub(time.Now()) < time.Hour*24*7 {
		_ = db.ExtendSessionLifetime(sessionID)
	}

	return session, nil
}

func (db *db.Database) DeleteSession(sessionID string) error {
	sessionIDByte := utils.Unhex(sessionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_, _ = db.db.ExecContext(
		ctx,
		"DELETE FROM sessions WHERE id = ?",
		sessionIDByte[:],
	)
	cancel()
	if ctx.Err() != nil && ctx.Err() != context.Canceled {
		log.Error(ctx.Err())
	}
	return nil
}

func (db *db.Database) DeleteExpiredSessions() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	_, err := db.db.ExecContext(
		ctx,
		"DELETE FROM sessions WHERE expires_at < ?",
		time.Now(),
	)
	cancel()
	if ctx.Err() != nil && ctx.Err() != context.Canceled {
		log.Error(ctx.Err())
	}
	return err
}

func (db *db.Database) ExtendSessionLifetime(sessionID string) error {
	sessionIDByte := utils.Unhex(sessionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_, err := db.db.ExecContext(
		ctx,
		"UPDATE sessions SET expires_at = ? WHERE id = ?",
		time.Now().Add(sessionDuration),
		sessionIDByte[:],
	)
	cancel()
	if ctx.Err() != nil {
		log.Error(ctx.Err())
	} else if err != nil {
		return ErrSessionExpired
	}
	return nil
}
