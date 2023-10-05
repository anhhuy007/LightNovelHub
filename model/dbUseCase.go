package model

type DB interface {
	CreateSession(userID []byte, deviceName string) (SessionInfo, bool)
	GetSession(sessionID []byte) (Session, bool)
	DeleteSession(sessionID []byte) bool
	DeleteExpiredSessions() bool
	DeleteAllSessions(userID []byte) bool
	ExtendSessionLifetime(sessionID []byte) bool

	CreateUser(username string, password []byte) ([]byte, bool)
	GetUser(username string) (User, bool)
	GetUserView(username string) (UserView, bool)
	GetUserViewByID(userID []byte) (UserView, bool)
	GetUserByID(userID []byte) (User, bool)
	GetUserMetadataSmall(userID []byte) (UserMetadataSmall, bool)
	FindUsers(username string, page uint) []UserMetadataSmall
	DeleteUser(userID []byte) bool
	UpdateUserMetadata(userID []byte, args *UserMetadata) bool
	UpdateUserPassword(userID []byte, newPassword []byte) bool

	GetFollowedUser(userID []byte) []UserMetadataSmall
	GetFollowedNovel(userID []byte, filtersAndSort *FiltersAndSortNovel) []NovelMetadataSmall

	CreateNovel(args *NovelMetadata) ([]byte, bool)
	GetNovelView(novelID []byte) (NovelView, bool)
	FindNovels(filtersAndSort *FiltersAndSortNovel) []NovelMetadataSmall
	UpdateNovelMetadata(novelID []byte, args *NovelMetadata) bool
	GetUsersNovels(
		userID []byte,
		filtersAndSort *FiltersAndSortNovel,
		isSelf bool,
	) []NovelMetadataSmall
}
