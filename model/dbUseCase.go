package model

type NovelMetadata struct {
	Title       string        `json:"title"`
	Tagline     string        `json:"tagline"`
	Description string        `json:"description"`
	Image       string        `json:"image"`
	Language    string        `json:"language"`
	Author      []byte        `json:"-"`
	Visibility  VisibilityID  `json:"visibility"`
	Status      NovelStatusID `json:"status"`
}

type DB interface {
	CreateSession(userID []byte, deviceName string) (SessionInfo, bool)
	GetSession(sessionID string) (Session, bool)
	DeleteSession(sessionID string) bool
	DeleteExpiredSessions() bool
	ExtendSessionLifetime(sessionID string) bool

	CreateUser(username, password string) ([]byte, bool)
	GetUser(username string) (User, bool)

	CreateNovel(args NovelMetadata) ([]byte, bool)
	GetNovelView(novelID string) (NovelView, bool)
	UpdateNovel(novelID string, args NovelMetadata) bool
}
