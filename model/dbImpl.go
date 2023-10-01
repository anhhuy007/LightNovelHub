package model

import (
	"Lightnovel/utils"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
	"time"
)

type Database struct {
	db              *sqlx.DB
	timeoutDuration time.Duration
}

func NewDatabase(db *sqlx.DB, timeoutDuration time.Duration) Database {
	return Database{
		db:              db,
		timeoutDuration: timeoutDuration,
	}
}

func (db *Database) GetUser(username string) (User, bool) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	err := db.db.GetContext(
		ctx,
		&user,
		"SELECT * FROM users WHERE username = ?",
		username,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return User{}, false
	}
	return user, true
}

func (db *Database) GetUserView(username string) (UserView, bool) {
	var userView UserView
	user, ok := db.GetUser(username)
	if !ok {
		return userView, false
	}
	userView.ID = user.ID
	userView.Username = user.Username
	userView.Image = user.Image
	userView.Displayname = user.Displayname
	userView.CreatedAt = user.CreatedAt

	return userView, true
}

func (db *Database) CreateUser(username, password string) ([]byte, bool) {
	var userId []byte

	hashed, err := utils.PasswordHash(password)
	if err != nil {
		log.Debug(err)
		return userId, false
	}

	userId = utils.GetUUID()
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	_, err = db.db.ExecContext(
		ctx,
		"INSERT INTO users (id, username, password) VALUES (?,?,?)",
		userId,
		username,
		hashed,
	)
	cancel()

	if err != nil && !errors.Is(err, context.Canceled) {
		log.Error(err)
		return userId, false
	}
	return userId, true
}

func (db *Database) CreateNovel(args NovelMetadata) ([]byte, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	uid := utils.GetUUID()
	_, err := db.db.ExecContext(
		ctx,
		"INSERT INTO novels (id, title, tagline, description, author, image, language, visibility) VALUES (?,?,?,?,?,?,?,?)",
		uid,
		args.Title,
		args.Tagline,
		args.Description,
		args.Author,
		args.Image,
		args.Language,
		args.Visibility,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return []byte{}, false
	}
	return uid, true
}

func (db *Database) getNovel(novelID []byte) (Novel, bool) {
	var novel Novel
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	err := db.db.GetContext(ctx, &novel, "SELECT * FROM novels WHERE id = ?", novelID)
	cancel()
	if err != nil {
		log.Error(err)
		return novel, false
	}
	return novel, true
}

func (db *Database) getAuthor(authorID []byte) (UserView, bool) {
	var user UserView
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	err := db.db.GetContext(
		ctx,
		&user,
		"SELECT id, username, displayname, image, created_at FROM users WHERE id = ?",
		authorID,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return user, false
	}
	return user, true
}

func (db *Database) getVisibility(visibilityID int) (string, bool) {
	var visibility string
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	err := db.db.GetContext(
		ctx,
		&visibility,
		"SELECT name FROM visibility WHERE id = ?",
		visibilityID,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return "", false
	}
	return visibility, true
}

func (db *Database) getStatus(statusId int) (string, bool) {
	var status string
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	err := db.db.GetContext(
		ctx,
		&status,
		"SELECT name FROM novel_status WHERE id = ?",
		statusId,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return "", false
	}
	return status, true
}

func (db *Database) getTags(novelID []byte) ([]TagView, bool) {
	var tags []TagView
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	row, err := db.db.QueryContext(
		ctx,
		"SELECT id, name FROM novel_tags LEFT JOIN tags on tags.id = novel_tags.tag_id WHERE novel_tags.novel_id = ?",
		novelID,
	)
	defer func() {
		err := row.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	for row.Next() {
		var tagView TagView
		err = row.Scan(&tagView.ID, &tagView.Name)
		if err != nil {
			log.Error(err)
			return tags, false
		}
		tags = append(tags, tagView)
	}
	cancel()
	return tags, true
}

func (db *Database) countVolume(novelID []byte) int {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	volumes := 0
	err := db.db.GetContext(
		ctx,
		&volumes,
		"SELECT COUNT(*) FROM volumes LEFT JOIN visibility ON volumes.visibility = visibility.id WHERE novel_id = ? AND visibility.name = 'PUB'",
		novelID,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return 0
	}
	return volumes
}

func (db *Database) GetNovelView(novelID string) (NovelView, bool) {
	novelIDBin := utils.Unhex(novelID)
	novel, ok := db.getNovel(novelIDBin)
	if !ok {
		return NovelView{}, false
	}

	author, ok := db.getAuthor(novel.Author)
	if !ok {
		return NovelView{}, false
	}

	tags, ok := db.getTags(novel.ID)
	if !ok {
		tags = []TagView{}
	}

	return NovelView{
		ID:          novelID,
		Title:       novel.Title,
		Tagline:     novel.Tagline,
		Description: novel.Description,
		Image:       novel.Image,
		Language:    novel.Language,
		CreateAt:    novel.CreateAt,
		UpdateAt:    novel.UpdateAt,
		TotalRating: novel.TotalRating,
		RateCount:   novel.RateCount,
		Views:       novel.Views,
		Clicks:      novel.Clicks,
		Adult:       novel.Adult,
		Author:      author,
		Status:      novel.StatusID.String(),
		Visibility:  novel.Visibility.String(),
		Tags:        tags,
		Volumes:     db.countVolume(novelIDBin),
	}, false
}

func (db *Database) UpdateNovel(novelID string, args NovelMetadata) bool {
	novelIDBin := utils.Unhex(novelID)
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	_, err := db.db.ExecContext(
		ctx,
		"UPDATE novels SET title = ?, tagline = ?, description = ?, image = ?, language = ?, visibility = ?, status_id = ? WHERE id = ?",
		args.Title,
		args.Tagline,
		args.Description,
		args.Image,
		args.Language,
		args.Visibility,
		args.Status,
		novelIDBin,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}
