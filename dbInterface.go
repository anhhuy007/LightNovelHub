package db

import "Lightnovel"

type DB interface {
	CreateSession(userID []byte, deviceName string) (main.SessionInfo, error)
	GetSession(sessionID string) (Session, error)
	DeleteSession(sessionID string) error
	DeleteExpiredSessions() error
	ExtendSessionLifetime(sessionID string) error

	CreateUser(username, password string) ([]byte, error)
	GetUser(username string) (User, error)
}
