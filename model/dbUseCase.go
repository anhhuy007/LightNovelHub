package model

type DB interface {
	CreateSession(userID []byte, deviceName string) (SessionInfo, bool)
	GetSession(sessionID []byte) (Session, bool)
	DeleteSession(sessionID []byte) bool
	DeleteExpiredSessions() bool
	DeletaAllSessions(userID []byte) bool
	ExtendSessionLifetime(sessionID []byte) bool

	CreateUser(username string, password []byte) ([]byte, bool)
	GetUser(username string) (User, bool)
	GetUserView(username string) (UserView, bool)
	GetUserViewWithID(userID []byte) (UserView, bool)
	GetUserWithID(userID []byte) (User, bool)
	GetUserMetadataSmall(userID []byte) (UserMetadataSmall, bool)
	DeleteUser(userID []byte) bool
	UpdateUserMetadata(userID []byte, args *UserMetadata) bool
	UpdateUserPassword(userID []byte, newPassword []byte) bool

	GetFollowedUser(userID []byte) []UserMetadataSmall
	GetFollowedNovel(userID []byte, filtersAndSort *FiltersAndSort) []NovelMetadataSmall

	CreateNovel(args *NovelMetadata) ([]byte, bool)
	GetNovelView(novelID []byte) (NovelView, bool)
	GetNovel(novelID []byte) (Novel, bool)
	UpdateNovelMetadata(novelID []byte, args *NovelMetadata) bool
	GetUsersNovels(UserID []byte, filtersAndSort *FiltersAndSort, isSelf bool) []NovelMetadataSmall
}
