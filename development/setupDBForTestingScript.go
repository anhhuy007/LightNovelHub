package main

import (
	"Lightnovel/model"
	"Lightnovel/utils"
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	mysqlConfig := mysql.Config{
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Addr:      os.Getenv("MYSQL_HOST"),
		Net:       "tcp",
		DBName:    os.Getenv("MYSQL_DATABASE"),
		ParseTime: true,
	}
	db, err := sqlx.ConnectContext(ctx, "mysql", mysqlConfig.FormatDSN())
	cancel()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if db.Ping() != nil {
		panic(err)
	}

	DropAllTable(db)
	SetupDatabase(db)

	// INSERT MOCK DATA
	users := []model.User{
		{
			ID:          utils.GetUUID(),
			Username:    "Thong",
			Displayname: "Thông Nguyễn",
			Password:    pwd("Thong12345"),
			Email:       sqlString("thong@thong.com"),
			Image:       "",
		},
		{
			ID:          utils.GetUUID(),
			Username:    "Huy",
			Displayname: "Anh Huy",
			Password:    pwd("AnhHuy12345"),
			Email:       sqlString("anh@huy.com"),
			Image:       "https://images.unsplash.com/photo-1693143600183-85bd2f7c99e5?auto=format&fit=crop&w=2070&q=80",
		},
		{
			ID:          utils.GetUUID(),
			Username:    "Vu",
			Displayname: "Anh Vũ",
			Password:    pwd("Guraa12345"),
			Email:       sqlString("vu@anh.com"),
			Image:       "https://images.unsplash.com/photo-1692260122105-28c26fc3c882?auto=format&fit=crop&w=720&q=80",
		},
	}
	for _, user := range users {
		db.MustExec(
			"INSERT INTO users (id, username, displayname, password, email, image) VALUES (?,?,?,?,?,?)",
			user.ID,
			user.Username,
			user.Displayname,
			user.Password,
			user.Email,
			user.Image,
		)
	}

	tags := []model.Tag{
		{
			ID:          1,
			Name:        "Action",
			Description: "Action",
		},
		{
			ID:          2,
			Name:        "Adventure",
			Description: "Adventure",
		},
		{
			ID:          3,
			Name:        "Comedy",
			Description: "Comedy",
		},
	}
	for _, tag := range tags {
		db.MustExec(
			"INSERT INTO tags (id, name, description) VALUES (?,?,?)",
			tag.ID, tag.Name, tag.Description,
		)
	}

	novels := []model.Novel{
		{
			ID:          utils.GetUUID(),
			Title:       "Tensei Shitara Slime Datta Ken",
			Tagline:     "That Time I Got Reincarnated as a Slime",
			Description: "Some description",
			Author:      users[0].ID,
			Image:       "https://images.unsplash.com/photo-1693346223929-17afbce70514?auto=format&fit=crop&w=1974&q=80",
			Language:    "eng",
			Visibility:  2,
			StatusID:    1,
			Adult:       false,
			Views:       50,
			Clicks:      70,
		},
		{
			ID:          utils.GetUUID(),
			Title:       "Solo Leveling",
			Tagline:     "Solo Leveling",
			Description: "Some description",
			Author:      users[1].ID,
			Image:       "https://images.unsplash.com/photo-1693369832705-3954d816601b?auto=format&fit=crop&w=1974&q=80",
			Language:    "eng",
			Visibility:  2,
			StatusID:    2,
			Adult:       false,
			Views:       100,
			Clicks:      200,
		},
		{
			ID:          utils.GetUUID(),
			Title:       "The Beginning After the End",
			Tagline:     "The Beginning After the End",
			Description: "Some description",
			Author:      users[2].ID,
			Image:       "",
			Language:    "eng",
			Visibility:  1,
			StatusID:    1,
			Adult:       true,
			Views:       40,
			Clicks:      50,
		},
		{
			ID:          utils.GetUUID(),
			Title:       "The Legendary Moonlight Sculptor",
			Tagline:     "The Legendary Moonlight Sculptor",
			Description: "Some description",
			Author:      users[2].ID,
			Image:       "https://images.unsplash.com/photo-1692893906137-9a93c125825c?auto=format&fit=crop&w=1964&q=80",
			Language:    "eng",
			Visibility:  2,
			StatusID:    1,
			Adult:       false,
			Views:       90,
			Clicks:      100,
		},
	}
	for _, novel := range novels {
		db.MustExec(
			"INSERT INTO novels (id, title, tagline, description, author, image, language, visibility, status_id, adult, views, clicks) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
			novel.ID,
			novel.Title,
			novel.Tagline,
			novel.Description,
			novel.Author,
			novel.Image,
			novel.Language,
			novel.Visibility,
			novel.StatusID,
			novel.Adult,
			novel.Views,
			novel.Clicks,
		)
	}

	novelTags := []model.Novel_tags{
		{
			NovelID: novels[0].ID,
			TagID:   tags[0].ID,
		},
		{
			NovelID: novels[0].ID,
			TagID:   tags[1].ID,
		},
		{
			NovelID: novels[1].ID,
			TagID:   tags[2].ID,
		},
		{
			NovelID: novels[2].ID,
			TagID:   tags[0].ID,
		},
	}
	for _, novelTag := range novelTags {
		db.MustExec(
			"INSERT INTO novel_tags (novel_id, tag_id) VALUES (?,?)",
			novelTag.NovelID, novelTag.TagID,
		)
	}

	volumes := []model.Volume{}
	for i := 0; i < 3; i++ {
		for j := 1; j <= 2; j++ {
			volumes = append(volumes, model.Volume{
				ID:          utils.GetUUID(),
				NovelID:     novels[i].ID,
				Title:       "Volume " + strconv.Itoa(j) + " of " + novels[i].Title,
				Tagline:     "Tagline",
				Description: "Description",
				Image:       novels[i].Image,
				Visibility:  2,
				Views:       50 + j,
			})
		}
		volumes = append(volumes, model.Volume{
			ID:          utils.GetUUID(),
			NovelID:     novels[i].ID,
			Title:       "Volume " + "3" + " of " + novels[i].Title,
			Tagline:     "Tagline",
			Description: "Description",
			Image:       novels[i].Image,
			Visibility:  1,
		})
	}
	for _, volume := range volumes {
		db.MustExec(
			"INSERT INTO volumes (id, novel_id, title, tagline, description, image, visibility, views) VALUES (?,?,?,?,?,?,?,?)",
			volume.ID,
			volume.NovelID,
			volume.Title,
			volume.Tagline,
			volume.Description,
			volume.Image,
			volume.Visibility,
			volume.Views,
		)
	}

	chapters := []model.Chapter{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			chapters = append(chapters, model.Chapter{
				ID:         utils.GetUUID(),
				VolumeID:   volumes[i*3+j].ID,
				Title:      "Chapter " + strconv.Itoa(j) + " of " + volumes[i*3+j].Title,
				Content:    "# Content\n ## Content\n ### Content\n #### Content",
				Visibility: 2,
				Views:      50 + j,
			})
		}
		chapters = append(chapters, model.Chapter{
			ID:         utils.GetUUID(),
			VolumeID:   volumes[i*3].ID,
			Title:      "Chapter " + "3" + " of " + volumes[i*3].Title,
			Content:    "# Content\n ## Content\n ### Content\n #### Content",
			Visibility: 1,
		})
	}
	for _, chapter := range chapters {
		db.MustExec(
			"INSERT INTO chapters (id, volume_id, title, content, visibility) VALUES (?,?,?,?,?)",
			chapter.ID, chapter.VolumeID, chapter.Title, chapter.Content, chapter.Visibility,
		)
	}

	images := []model.Image{
		{
			ID:      1,
			UserID:  users[0].ID,
			NovelID: novels[0].ID,
			Url:     "https://plus.unsplash.com/premium_photo-1692878409491-e20dca376c97?auto=format&fit=crop&w=1974&q=80",
		},
		{
			ID:      2,
			UserID:  users[1].ID,
			NovelID: novels[1].ID,
			Url:     "https://images.unsplash.com/photo-1693322565356-eaa6771e2442?auto=format&fit=crop&w=1974&q=80",
		},
		{
			ID:      3,
			UserID:  users[2].ID,
			NovelID: novels[2].ID,
			Url:     "https://images.unsplash.com/photo-1693275449979-aa6d53df26b8?auto=format&fit=crop&w=1974&q=80",
		},
	}
	for _, image := range images {
		db.MustExec(
			"INSERT INTO images (id, user_id, novel_id, url) VALUES (?,?,?,?)",
			image.ID, image.UserID, image.NovelID, image.Url,
		)
	}

	comments := []model.Comment{}
	for _, user := range users {
		for _, novel := range novels {
			id := utils.GetUUID()
			comments = append(comments, model.Comment{
				ID:      id,
				ToID:    novel.ID,
				UserID:  user.ID,
				Content: "**Comment** content",
			})
			comments = append(comments, model.Comment{
				ID:      utils.GetUUID(),
				ToID:    id,
				UserID:  user.ID,
				Content: "Commenting to a *comment*",
			})
		}
		for _, volume := range volumes {
			comments = append(comments, model.Comment{
				ID:      utils.GetUUID(),
				ToID:    volume.ID,
				UserID:  user.ID,
				Content: "Comment *content*",
			})
		}
		for _, chapter := range chapters {
			comments = append(comments, model.Comment{
				ID:      utils.GetUUID(),
				ToID:    chapter.ID,
				UserID:  user.ID,
				Content: "Comment content",
			})
		}
	}
	for _, comment := range comments {
		db.MustExec(
			"INSERT INTO comments (id, to_id, user_id, content) VALUES (?,?,?,?)",
			comment.ID, comment.ToID, comment.UserID, comment.Content,
		)
	}

}

func DropAllTable(db *sqlx.DB) {
	rows, _ := db.Query(
		"SELECT concat('DROP TABLE IF EXISTS `', table_name, '`;') AS result FROM information_schema.tables WHERE table_schema = ?;",
		os.Getenv("MYSQL_DATABASE"),
	)
	defer rows.Close()
	db.MustExec("SET FOREIGN_KEY_CHECKS = 0")
	for rows.Next() {
		s := ""
		_ = rows.Scan(&s)
		db.MustExec(s)
	}
	db.MustExec("SET FOREIGN_KEY_CHECKS = 1")
}

func ProcessRawSql(sql string) []string {
	tmp := strings.Split(sql, ";")
	res := make([]string, 0, len(tmp))
	for _, item := range tmp {
		replaced := strings.ReplaceAll(item, "\n", "")
		if len(replaced) != 0 {
			res = append(res, replaced)
		}
	}
	return res
}

func SetupDatabase(db *sqlx.DB) {
	startupScript, err := os.ReadFile("../dbScript/Schema.sql")
	if err != nil {
		panic(err)
	}
	lines := ProcessRawSql(string(startupScript))
	for _, line := range lines {
		db.MustExec(line)
	}
}

func pwd(s string) []byte {
	res, _ := utils.PasswordHash(s)
	return res
}

func sqlString(s string) sql.NullString {
	return sql.NullString{
		Valid:  true,
		String: s,
	}
}
