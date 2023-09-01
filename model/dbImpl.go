package model

import (
	"Lightnovel/utils"
	"database/sql"
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

func isLetterOnly(s string) bool {
	for _, c := range s {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
			return false
		}
	}
	return true
}

func isUsernameValid(username string) bool {
	return len(username) != 0 && len(username) <= 32 && isLetterOnly(username)
}

func isPasswordValid(password string) bool {
	return !(len(password) < 10 || len(password) > 72)
}

func (db *Database) GetUser(username string) (User, error) {
	var user User
	if !isUsernameValid(username) {
		return user, ErrInvalidUsername
	}
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	err := db.db.GetContext(
		ctx,
		&user,
		"SELECT * FROM users WHERE username = ?",
		username,
	)
	cancel()
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	} else if err != nil && !errors.Is(err, context.Canceled) {
		log.Error(err)
		return User{}, ErrInternal
	}
	return user, nil
}

func (db *Database) GetUserView(username string) (UserView, error) {
	var userView UserView
	user, err := db.GetUser(username)
	if err != nil {
		return userView, err
	}
	userView.ID = user.ID
	userView.Username = user.Username
	userView.Image = user.Image
	userView.Displayname = user.Displayname
	userView.CreatedAt = user.CreatedAt

	return userView, nil
}

func (db *Database) CreateUser(username, password string) ([]byte, error) {
	var userId []byte
	_, err := db.GetUser(username)
	if err == nil {
		return userId, ErrUserAlreadyExist
	}

	if !isPasswordValid(password) {
		return userId, ErrInvalidPassword
	}
	hashed, err := utils.PasswordHash(password)
	if err != nil {
		log.Debug(err)
		return userId, ErrInvalidPassword
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
		return userId, ErrInternal
	}
	return userId, nil
}

func (db *Database) CreateNovel(args CreateNovelArgs) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	uid := utils.GetUUID()
	_, err := db.db.ExecContext(
		ctx,
		"INSERT INTO novels (id, title, tagline, description, author, image, language) VALUES (?,?,?,?,?,?,?)",
		uid,
		args.Title,
		args.Tagline,
		args.Description,
		args.Author,
		args.Image,
		args.Language,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return []byte{}, ErrInternal
	}
	return uid, nil
}

func (db *Database) getNovel(novelID []byte) (Novel, error) {
	var novel Novel
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	err := db.db.GetContext(ctx, &novel, "SELECT * FROM novels WHERE id = ?", novelID)
	cancel()
	if err != nil {
		return novel, err
	}
	return novel, nil
}

func (db *Database) getAuthor(authorID []byte) (UserView, error) {
	var user UserView
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	err := db.db.GetContext(
		ctx,
		&user,
		"SELECT id, username, displayname, image, created_at FROM users WHERE id = ?",
		authorID,
	)
	cancel()
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrUserNotFound
	} else if err != nil && !errors.Is(err, context.Canceled) {
		log.Error(err)
		return user, ErrInternal
	}
	return user, nil
}

func (db *Database) getVisibility(visibilityID int) (string, error) {
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
		return "", ErrInternal
	}
	return visibility, nil
}

func (db *Database) getStatus(statusId int) (string, error) {
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
		return "", ErrInternal
	}
	return status, nil
}

func (db *Database) getTags(novelID []byte) ([]TagView, error) {
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
			return tags, ErrInternal
		}
		tags = append(tags, tagView)
	}
	cancel()
	return tags, nil
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
	}
	return volumes
}

func (db *Database) GetNovelView(novelID string) (NovelView, error) {
	novelIDBin := utils.Unhex(novelID)
	novel, err := db.getNovel(novelIDBin)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error(err)
			return NovelView{}, ErrNovelNotFound
		}
		return NovelView{}, ErrInternal
	}

	author, err := db.getAuthor(novel.Author)
	if err != nil {
		return NovelView{}, err
	}

	status, err := db.getStatus(novel.StatusID)
	if err != nil {
		return NovelView{}, err
	}
	visibility, err := db.getVisibility(novel.Visibility)
	if err != nil {
		return NovelView{}, err
	}

	tags, err := db.getTags(novel.ID)
	if err != nil {
		return NovelView{}, err
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
		Status:      status,
		Visibility:  visibility,
		Tags:        tags,
		Volumes:     db.countVolume(novelIDBin),
	}, nil
}
