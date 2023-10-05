package repo

import (
	"Lightnovel/model"
	"context"
	"encoding/hex"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

var sessionDuration = time.Hour * 24 * 30

func (db *Database) CreateSession(
	userID []byte,
	deviceName string,
) (model.SessionInfo, bool) {
	sessionID := GetUUID()
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
		return model.SessionInfo{}, false
	}
	return model.SessionInfo{hex.EncodeToString(sessionID), expires}, true
}

func (db *Database) GetSession(sessionID []byte) (model.Session, bool) {
	var session model.Session
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	err := db.db.GetContext(
		ctx,
		&session,
		"SELECT id, user_id, expires_at, device_name FROM sessions WHERE id = ?",
		sessionID,
	)
	cancel()

	if err != nil {
		return model.Session{}, false
	}

	if session.ExpireAt.Sub(time.Now()) < sessionDuration/3 {
		_ = db.ExtendSessionLifetime(sessionID)
	}

	return session, true
}

func (db *Database) DeleteSession(sessionID []byte) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_, err := db.db.ExecContext(
		ctx,
		"DELETE FROM sessions WHERE id = ?",
		sessionID,
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

func (db *Database) ExtendSessionLifetime(sessionID []byte) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_, err := db.db.ExecContext(
		ctx,
		"UPDATE sessions SET expires_at = ? WHERE id = ?",
		time.Now().Add(sessionDuration),
		sessionID,
	)
	cancel()
	if err != nil {
		log.Error(ctx.Err())
		return false
	}
	return true
}

func (db *Database) DeleteAllSessions(userID []byte) bool {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	_, err := db.db.ExecContext(ctx, "DELETE FROM sessions WHERE user_id = ?", userID)
	cancel()
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}
