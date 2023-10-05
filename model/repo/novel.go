package repo

import (
	"Lightnovel/model"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
)

func (db *Database) CreateNovel(args *model.NovelMetadata) ([]byte, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	uid := GetUUID()
	_, err := db.db.ExecContext(
		ctx,
		`INSERT INTO novels 
        (id, title, tagline, description, author, image, language, visibility) 
		VALUES (?,?,?,?,?,?,?,?)`,
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

func (db *Database) GetNovel(novelID []byte) (model.Novel, bool) {
	var novel model.Novel
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	err := db.db.GetContext(ctx, &novel, "SELECT * FROM novels WHERE id = ?", novelID)
	cancel()
	if err != nil {
		log.Error(err)
		return novel, false
	}
	return novel, true
}

func (db *Database) getAuthor(authorID []byte) (model.UserView, bool) {
	var user model.UserView
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	err := db.db.GetContext(
		ctx,
		&user,
		`SELECT id, username, displayname, image, created_at 
		FROM users 
		WHERE id = ?`,
		authorID,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return user, false
	}
	return user, true
}

func (db *Database) getTags(novelID []byte) []model.TagView {
	var tags []model.TagView
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	row, err := db.db.QueryxContext(
		ctx,
		`SELECT id, name 
		FROM novel_tags LEFT JOIN tags 
		ON tags.id = novel_tags.tag_id 
		WHERE novel_tags.novel_id = ?`,
		novelID,
	)
	defer func() {
		err := row.Close()
		if err != nil {
			log.Error(err)
		}
		cancel()
	}()
	for row.Next() {
		var tagView model.TagView
		err = row.StructScan(&tagView)
		if err != nil {
			log.Error(err)
			return tags
		}
		tags = append(tags, tagView)
	}
	return tags
}

func (db *Database) countVolume(novelID []byte) int {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	volumes := 0
	err := db.db.GetContext(
		ctx,
		&volumes,
		`SELECT COUNT(*) 
		FROM volumes
		WHERE novel_id = ? AND visibility = ?`,
		novelID,
		model.VisibilityPublic,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return 0
	}
	return volumes
}

func (db *Database) GetNovelView(novelID []byte) (model.NovelView, bool) {
	novel, ok := db.GetNovel(novelID)
	if !ok {
		return model.NovelView{}, false
	}

	author, ok := db.GetUserMetadataSmall(novel.Author)
	if !ok {
		return model.NovelView{}, false
	}

	return model.NovelView{
		ID:          hex.EncodeToString(novelID),
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
		Status:      novel.Status.String(),
		Visibility:  novel.Visibility.String(),
		Tags:        db.getTags(novelID),
		Volumes:     db.countVolume(novelID),
		FollowCount: db.countUserFollowers(novelID),
	}, true
}

func (db *Database) UpdateNovelMetadata(novelID []byte, args *model.NovelMetadata) bool {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	_, err := db.db.ExecContext(
		ctx,
		`UPDATE novels 
		SET title = ?, tagline = ?, description = ?, image = ?, 
		    language = ?, visibility = ?, status = ? 
		WHERE id = ?`,
		args.Title,
		args.Tagline,
		args.Description,
		args.Image,
		args.Language,
		args.Visibility,
		args.Status,
		novelID,
	)
	cancel()
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (db *Database) GetUsersNovels(
	userID []byte,
	filtersAndSort *model.FiltersAndSortNovel,
	isSelf bool,
) []model.NovelMetadataSmall {
	var novels []model.NovelMetadataSmall
	filtersAndSortQuery, filtersAndSortArgs := filtersAndSort.ConstructQuery()
	query := `
		SELECT novels.* 
		FROM novels
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
	query += ` WHERE novels.author = ?`
	if isSelf == false {
		query += fmt.Sprintf(" AND novels.visibility = %v", int(model.VisibilityPublic))
	}
	query += filtersAndSortQuery

	args := []interface{}{userID}
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

func (db *Database) FindNovels(
	filtersAndSort *model.FiltersAndSortNovel,
) []model.NovelMetadataSmall {
	var novels []model.NovelMetadataSmall
	filtersAndSortQuery, filtersAndSortArgs := filtersAndSort.ConstructQuery()
	query := `
		SELECT novels.* 
		FROM novels
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
	query += " WHERE 1=1 "
	query += filtersAndSortQuery

	ctx, cancel := context.WithTimeout(context.Background(), db.timeoutDuration)
	row, err := db.db.QueryxContext(
		ctx,
		query,
		filtersAndSortArgs...,
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
