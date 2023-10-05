package repo

import (
	"Lightnovel/model"
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

func (db *Database) GetUser(username string) (model.User, bool) {
	var user model.User
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
		return model.User{}, false
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

func (db *Database) GetUserView(username string) (model.UserView, bool) {
	user, ok := db.GetUser(username)
	if !ok {
		return model.UserView{}, false
	}
	userView := model.UserView{
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

func (db *Database) GetUserViewByID(userID []byte) (model.UserView, bool) {
	user, ok := db.GetUserByID(userID)
	if !ok {
		return model.UserView{}, false
	}
	userView := model.UserView{
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

func (db *Database) GetUserByID(userID []byte) (model.User, bool) {
	var user model.User
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	err := db.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", userID)
	cancel()
	if err != nil {
		log.Error(err)
		return user, false
	}
	return user, true
}

func (db *Database) GetUserMetadataSmall(userID []byte) (model.UserMetadataSmall, bool) {
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
		return model.UserMetadataSmall{}, false
	}
	return model.UserMetadataSmall{
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

func (db *Database) UpdateUserMetadata(userID []byte, args *model.UserMetadata) bool {
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

type UserMetadatSmallRaw struct {
	Id          []byte
	Username    string
	Displayname sql.NullString
	Image       string
}

func (db *Database) GetFollowedUser(userID []byte) []model.UserMetadataSmall {
	var users []model.UserMetadataSmall
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	row, err := db.db.QueryxContext(
		ctx,
		`SELECT 
    			users.id, users.username, users.displayname, users.image 
		FROM follows_user LEFT JOIN users 
		ON follows_user.to_id = users.id 
		WHERE from_id = ? 
		ORDER BY users.username`,
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
		var userMetaSmallRaw UserMetadatSmallRaw
		err := row.StructScan(&userMetaSmallRaw)
		if err != nil {
			log.Error(err)
			return users
		}
		users = append(users, model.UserMetadataSmall{
			ID:          hex.EncodeToString(userMetaSmallRaw.Id),
			Username:    userMetaSmallRaw.Username,
			Displayname: userMetaSmallRaw.Displayname.String,
			Image:       userMetaSmallRaw.Image,
		})
	}
	return users
}

func (db *Database) GetFollowedNovel(
	userID []byte,
	filtersAndSort *model.FiltersAndSortNovel,
) []model.NovelMetadataSmall {
	var novels []model.NovelMetadataSmall
	filtersAndSortQuery, filtersAndSortArgs := filtersAndSort.ConstructQuery()
	query := `
		SELECT novels.* 
		FROM follows_novel 
		LEFT JOIN novels 
		ON follows_novel.novel_id = novels.id
	`
	if len(filtersAndSort.Tag) != 0 || len(filtersAndSort.TagExclude) != 0 {
		query += `
		RIGHT JOIN (
			SELECT novel_id, GROUP_CONCAT(tag_id) AS tag_groupconcat
			FROM novel_tags
			GROUP BY novel_id
		) AS TABLE1
		ON TABLE1.novel_id = novels.id`
	}
	query += ` WHERE follows_novel.user_id = ? AND novels.visibility = ?` + filtersAndSortQuery

	args := []interface{}{userID, model.VisibilityPublic}
	if filtersAndSortArgs != nil {
		args = append(args, filtersAndSortArgs...)
	}

	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	row, err := db.db.QueryxContext(
		ctx,
		query,
		args...,
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
		var novel model.Novel
		err := row.StructScan(&novel)
		if err != nil {
			log.Error(err)
			return novels
		}
		authorMetadataSmall, ok := db.GetUserMetadataSmall(novel.Author)
		if !ok {
			return novels
		}
		novels = append(novels, model.NovelMetadataSmall{
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

func (db *Database) FindUsers(username string, page uint) []model.UserMetadataSmall {
	var usersMetadataSmall []model.UserMetadataSmall
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	rows, err := db.db.QueryxContext(
		ctx,
		`SELECT id, username, displayname, image FROM users WHERE username LIKE ? ORDER BY username LIMIT ? OFFSET ?`,
		username,
		model.PageSize,
		model.PageSize*(page-1),
	)
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error(err)
		}
		cancel()
	}()
	if err != nil {
		log.Error(err)
		return usersMetadataSmall
	}
	for rows.Next() {
		var raw UserMetadatSmallRaw
		if err := rows.StructScan(&raw); err != nil {
			log.Error(err)
			return usersMetadataSmall
		}
		usersMetadataSmall = append(usersMetadataSmall, model.UserMetadataSmall{
			ID:          hex.EncodeToString(raw.Id),
			Username:    raw.Username,
			Displayname: raw.Displayname.String,
			Image:       raw.Image,
		})
	}

	return usersMetadataSmall
}
