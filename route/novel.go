package route

import (
	"Lightnovel/middleware"
	"Lightnovel/model"
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func AddUploadRoutes(router *fiber.Router, db model.DB) {
	novelRoute := (*router).Group("/novel")

	type CreateNovelResult struct {
		NovelID string `json:"novel_id"`
	}
	novelRoute.Post("/create", func(c *fiber.Ctx) error {
		if c.Locals(middleware.KeyIsUserAuth) == false {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		var input model.CreateNovelArgs
		err := c.BodyParser(&input)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		session, ok := c.Locals(middleware.KeyUserSession).(model.Session)
		if !ok {
			log.Warn("Check the authentication middleware")
		}
		input.Author = session.UserID
		uid, err := db.CreateNovel(input)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusCreated).JSON(
			CreateNovelResult{
				NovelID: hex.EncodeToString(uid),
			})
	})

	novelRoute.Get("/:novelID", func(c *fiber.Ctx) error {
		novelID := c.Params("novelID")
		if len(novelID) != 32 {
			return c.SendStatus(fiber.StatusNotFound)
		}

		novelView, err := db.GetNovelView(novelID)

		if errors.Is(err, model.ErrNovelNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if novelView.Visibility == "PRI" {
			if c.Locals(middleware.KeyIsUserAuth) == false {
				return c.SendStatus(fiber.StatusUnauthorized)
			}
			sessionInfo, ok := c.Locals(middleware.KeyUserSession).(model.SessionInfo)
			if !ok {
				log.Warn("Check auth middleware")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			session, err := db.GetSession(sessionInfo.Session)
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			if bytes.Compare(session.UserID, novelView.Author.ID) != 0 {
				return c.SendStatus(fiber.StatusUnauthorized)
			}
		}

		// Everything is good
		return c.JSON(novelView)
	})
}
