package model

type CreateNovelArgs struct {
	Title       string
	Tagline     string
	Description string
	Image       string
	Language    string
	Author      []byte
}

type DB interface {
	CreateSession(userID []byte, deviceName string) (SessionInfo, error)
	GetSession(sessionID string) (Session, error)
	DeleteSession(sessionID string) error
	DeleteExpiredSessions() error
	ExtendSessionLifetime(sessionID string) error

	CreateUser(username, password string) ([]byte, error)
	GetUser(username string) (User, error)

	CreateNovel(args CreateNovelArgs) ([]byte, error)
	GetNovelView(novelID string) (NovelView, error)
}
