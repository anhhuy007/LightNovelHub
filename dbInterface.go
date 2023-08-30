package main

type DB interface {
	CreateSession(userID [16]byte, deviceName string) (string, error)
	GetSession(sessionID string) (Session, error)
	DeleteSession(sessionID string) error
	DeleteExpiredSessions() error
	ExtendSessionLifetime(sessionID string) error

	CreateUser(username, password string) error
	GetUser(username string) (User, error)
}
