package model

import (
	"Lightnovel/utils"
	"context"
	"database/sql"
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
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Error(err)
		return SessionInfo{}, ErrInternal
	}
	return SessionInfo{hex.EncodeToString(sessionID), expires}, nil
}

func (db *Database) GetSession(sessionID string) (Session, error) {
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

	if errors.Is(err, sql.ErrNoRows) {
		return Session{}, ErrSessionExpired
	} else if err != nil && !errors.Is(err, context.Canceled) {
		log.Error(err)
		return Session{}, ErrInternal
	}
	if time.Now().After(session.ExpireAt) {
		_ = db.DeleteSession(sessionID)
		return Session{}, ErrSessionExpired
	}

	if session.ExpireAt.Sub(time.Now()) < time.Hour*24*7 {
		_ = db.ExtendSessionLifetime(sessionID)
	}

	return session, nil
}

func (db *Database) DeleteSession(sessionID string) error {
	sessionIDByte := utils.Unhex(sessionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_, _ = db.db.ExecContext(
		ctx,
		"DELETE FROM sessions WHERE id = ?",
		sessionIDByte,
	)
	cancel()
	if ctx.Err() != nil && !errors.Is(ctx.Err(), context.Canceled) {
		log.Error(ctx.Err())
	}
	return nil
}

func (db *Database) DeleteExpiredSessions() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	_, err := db.db.ExecContext(
		ctx,
		"DELETE FROM sessions WHERE expires_at < ?",
		time.Now(),
	)
	cancel()
	if ctx.Err() != nil && !errors.Is(ctx.Err(), context.Canceled) {
		log.Error(ctx.Err())
	}
	return err
}

func (db *Database) ExtendSessionLifetime(sessionID string) error {
	sessionIDByte := utils.Unhex(sessionID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_, err := db.db.ExecContext(
		ctx,
		"UPDATE sessions SET expires_at = ? WHERE id = ?",
		time.Now().Add(sessionDuration),
		sessionIDByte,
	)
	cancel()
	if ctx.Err() != nil {
		log.Error(ctx.Err())
	} else if err != nil {
		return ErrSessionExpired
	}
	return nil
}
