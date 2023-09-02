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

type createNovelResult struct {
	NovelID string `json:"novel_id"`
}

func AddUploadRoutes(router *fiber.Router, db model.DB) {
	novelRoute := (*router).Group("/novel")

	novelRoute.Post("/create", createNovel(db))

	novelRoute.Get("/:novelID", getNovel(db))
}

// Create Novel
//
//	@Summary	Create a new novel, return the created novel id
//	@Tags		novel
//	@Accept		json
//	@Produce	json
//	@Param		NovelDetails	body		model.CreateNovelArgs	true	"Novel details"
//	@Success	200				{object}	createNovelResult
//	@Failure	400
//	@Failure	401
//	@Failure	500
//	@Router		/novel/create [POST]
func createNovel(db model.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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
			createNovelResult{
				NovelID: hex.EncodeToString(uid),
			})
	}
}

// Get Novel
//
//	@Summary		Get the novel with provided novel id
//	@Description	If the novel is private, the user need to be logged in with the author account
//	@Tags			novel
//	@Produce		json
//	@Param			NovelID	path		string	true	"Novel ID"
//	@Success		200		{object}	model.NovelView
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/novel/:NovelID [GET]
func getNovel(db model.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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
	}
}
