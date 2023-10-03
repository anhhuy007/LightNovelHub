package model

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/gofiber/fiber/v2/log"
)

func (db *Database) CreateUser(username string, password []byte) ([]byte, bool) {
	userId := GetUUID()
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	_, err := db.db.ExecContext(
		ctx,
		"INSERT INTO users (id, username, password) VALUES (?,?,?)",
		userId,
		username,
		password,
	)
	cancel()

	if err != nil && !errors.Is(err, context.Canceled) {
		log.Error(err)
		return userId, false
	}
	return userId, true
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

func (db *Database) countUserNovel(userID []byte) int {
	novelCount := 0
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	err := db.db.GetContext(
		ctx,
		&novelCount,
		"SELECT COUNT(*) FROM novels WHERE author = ?",
		userID,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return 0
	}
	return novelCount
}

func (db *Database) GetUserView(username string) (UserView, bool) {
	user, ok := db.GetUser(username)
	if !ok {
		return UserView{}, false
	}
	userView := UserView{
		ID:            hex.EncodeToString(user.ID),
		Username:      user.Username,
		Displayname:   user.Displayname.String,
		Image:         user.Image,
		CreatedAt:     user.CreatedAt,
		NovelCount:    db.countUserNovel(user.ID),
		FollowerCount: db.countUserFollowers(user.ID),
		FollowedCount: db.countUserFollows(user.ID),
	}

	return userView, true
}

func (db *Database) GetUserViewWithID(userID []byte) (UserView, bool) {
	user, ok := db.GetUserWithID(userID)
	if !ok {
		return UserView{}, false
	}
	userView := UserView{
		ID:            hex.EncodeToString(user.ID),
		Username:      user.Username,
		Displayname:   user.Displayname.String,
		Image:         user.Image,
		CreatedAt:     user.CreatedAt,
		NovelCount:    db.countUserNovel(user.ID),
		FollowerCount: db.countUserFollowers(user.ID),
		FollowedCount: db.countUserFollows(user.ID),
	}

	return userView, true
}

func (db *Database) GetUserWithID(userID []byte) (User, bool) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	err := db.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", userID)
	cancel()
	if err != nil {
		log.Error(err)
		return user, false
	}
	return user, true
}

func (db *Database) GetUserMetadataSmall(userID []byte) (UserMetadataSmall, bool) {
	var userMetadataSmall struct {
		ID          []byte
		Username    string
		Displayname sql.NullString
		Image       string
	}
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	err := db.db.GetContext(
		ctx,
		&userMetadataSmall,
		"SELECT id, username, displayname, image FROM users WHERE id = ?",
		userID,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return UserMetadataSmall{}, false
	}
	return UserMetadataSmall{
		ID:          hex.EncodeToString(userMetadataSmall.ID),
		Username:    userMetadataSmall.Username,
		Displayname: userMetadataSmall.Displayname.String,
		Image:       userMetadataSmall.Image,
	}, true
}

func (db *Database) DeleteUser(userID []byte) bool {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	_, err := db.db.ExecContext(
		ctx,
		"DELETE FROM users WHERE id = ?",
		userID,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (db *Database) UpdateUserMetadata(userID []byte, args UserMetadata) bool {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	_, err := db.db.ExecContext(
		ctx,
		"UPDATE users SET username = ?, displayname = ?, email = ?, image = ? WHERE id = ?",
		args.Username,
		args.Displayname,
		args.Email,
		args.Image,
		userID,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (db *Database) UpdateUserPassword(userID []byte, newPassword []byte) bool {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	_, err := db.db.ExecContext(
		ctx,
		"UPDATE users SET password = ? WHERE id = ?",
		newPassword,
		userID,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (db *Database) GetFollowedUser(userID []byte) []UserMetadataSmall {
	var users []UserMetadataSmall
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	row, err := db.db.QueryxContext(
		ctx,
		"SELECT users.id, users.username, users.displayname, users.image FROM follows_user LEFT JOIN users ON follows_user.to_id = users.id WHERE from_id = ?",
		userID,
	)
	defer func() {
		err := row.Close()
		if err != nil {
			log.Error(err)
		}
		cancel()
	}()
	if err != nil {
		log.Error(err)
		return users
	}
	for row.Next() {
		var userMetaSmallRaw struct {
			Id          []byte
			Username    string
			Displayname sql.NullString
			Image       string
		}
		err := row.StructScan(&userMetaSmallRaw)
		if err != nil {
			log.Error(err)
			return users
		}
		users = append(users, UserMetadataSmall{
			ID:          hex.EncodeToString(userMetaSmallRaw.Id),
			Username:    userMetaSmallRaw.Username,
			Displayname: userMetaSmallRaw.Displayname.String,
			Image:       userMetaSmallRaw.Image,
		})
	}
	return users
}

func (db *Database) GetFollowedNovel(userID []byte) []NovelMetadataSmall {
	var novels []NovelMetadataSmall
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	row, err := db.db.QueryxContext(
		ctx,
		"SELECT novels.* FROM follows_novel LEFT JOIN novels ON follows_novel.novel_id = novels.id WHERE user_id = ? AND visibility = ?",
		userID,
		VisibilityPublic,
	)
	defer func() {
		err := row.Close()
		if err != nil {
			log.Error(err)
		}
		cancel()
	}()
	if err != nil {
		log.Error(err)
		return novels
	}
	for row.Next() {
		var novel Novel
		err := row.StructScan(&novel)
		if err != nil {
			log.Error(err)
			return novels
		}
		authorMetadataSmall, ok := db.GetUserMetadataSmall(userID)
		if !ok {
			return novels
		}
		novels = append(novels, NovelMetadataSmall{
			ID:          hex.EncodeToString(novel.ID),
			Title:       novel.Title,
			Tagline:     novel.Tagline,
			Description: novel.Description,
			Author:      authorMetadataSmall,
			Image:       novel.Image,
			Language:    novel.Language,
			TotalRating: novel.TotalRating,
			RateCount:   novel.RateCount,
			Adult:       novel.Adult,
			Status:      novel.Status.String(),
			Visibility:  novel.Visibility.String(),
			Views:       novel.Views,
		})
	}
	return novels
}
