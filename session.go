package main

import (
	"context"
	"encoding/hex"
	"time"
)

var sessionDuration = time.Hour * 24 * 30

func (db *Database) CreateSession(userID [16]byte, deviceName string) (string, error) {
	sessionID := getUUID()
	var tmp context.Context
	ctx, cancel := context.WithTimeout(tmp, time.Minute)
	_, err := db.db.ExecContext(ctx, "INSERT INTO sessions (id, user_id, expires_at, device_name) VALUES (?, ?, ?, ?)", sessionID, userID, time.Now().Add(sessionDuration), deviceName)
	cancel()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sessionID[:]), nil
}

func (db *Database) GetSession(sessionID string) (Session, error) {
	var session Session
	var tmp context.Context
	ctx, cancel := context.WithTimeout(tmp, time.Minute)
	err := db.db.GetContext(ctx, &session, "SELECT id, user_id, expires_at, device_name FROM sessions WHERE id = ?", unhex(sessionID))
	cancel()
	return session, err
}

func (db *Database) DeleteSession(sessionID string) error {
	var tmp context.Context
	ctx, cancel := context.WithTimeout(tmp, time.Minute)
	_, err := db.db.ExecContext(ctx, "DELETE FROM sessions WHERE id = ?", unhex(sessionID))
	cancel()
	return err
}

func (db *Database) DeleteExpiredSessions() error {
	var tmp context.Context
	ctx, cancel := context.WithTimeout(tmp, time.Minute*10)
	_, err := db.db.ExecContext(ctx, "DELETE FROM sessions WHERE expires_at < ?", time.Now())
	cancel()
	return err
}

func (db *Database) ExtendSessionLifetime(sessionID string) error {
	var tmp context.Context
	ctx, cancel := context.WithTimeout(tmp, time.Minute)
	_, err := db.db.ExecContext(ctx, "UPDATE sessions SET expires_at = ? WHERE id = ?", time.Now().Add(sessionDuration), unhex(sessionID))
	cancel()
	return err
}
