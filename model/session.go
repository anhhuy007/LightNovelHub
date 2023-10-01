package model

import (
	"Lightnovel/utils"
	"context"
	"encoding/hex"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

var sessionDuration = time.Hour * 24 * 30

type SessionInfo struct {
	Session   string    `json:"session"`
	ExpiredAt time.Time `json:"expired_at"`
}

type IncludeSessionString struct {
	Session string `json:"session"`
}

func (db *Database) CreateSession(
	userID []byte,
	deviceName string,
) (SessionInfo, bool) {
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
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Error(err)
		return SessionInfo{}, false
	}
	return SessionInfo{hex.EncodeToString(sessionID), expires}, true
}

func (db *Database) GetSession(sessionID string) (Session, bool) {
	var session Session
	sessionIDByte := utils.Unhex(sessionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	err := db.db.GetContext(
		ctx,
		&session,
		"SELECT id, user_id, expires_at, device_name FROM sessions WHERE id = ?",
		sessionIDByte,
	)
	cancel()

	if err != nil {
		return Session{}, false
	}

	if session.ExpireAt.Sub(time.Now()) < sessionDuration/3 {
		_ = db.ExtendSessionLifetime(sessionID)
	}

	return session, true
}

func (db *Database) DeleteSession(sessionID string) bool {
	sessionIDByte := utils.Unhex(sessionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_, err := db.db.ExecContext(
		ctx,
		"DELETE FROM sessions WHERE id = ?",
		sessionIDByte,
	)
	cancel()
	if err != nil {
		log.Error(ctx.Err())
	}
	return true
}

func (db *Database) DeleteExpiredSessions() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	_, err := db.db.ExecContext(
		ctx,
		"DELETE FROM sessions WHERE expires_at < ?",
		time.Now(),
	)
	cancel()
	if err != nil {
		log.Error(ctx.Err())
	}
	return true
}

func (db *Database) ExtendSessionLifetime(sessionID string) bool {
	sessionIDByte := utils.Unhex(sessionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_, err := db.db.ExecContext(
		ctx,
		"UPDATE sessions SET expires_at = ? WHERE id = ?",
		time.Now().Add(sessionDuration),
		sessionIDByte,
	)
	cancel()
	if err != nil {
		log.Error(ctx.Err())
		return false
	}
	return true
}
