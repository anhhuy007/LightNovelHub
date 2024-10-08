package route

import (
	"Lightnovel/middleware"
	"Lightnovel/model"
	"bytes"
	"encoding/hex"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"regexp"
	"unicode/utf8"
)

// TODO: Delete novel
func AddUploadRoutes(router *fiber.Router, db model.DB) {
	novelRoute := (*router).Group("/novel")

	novelRoute.Get("/find", searchAndFilterNovel(db))

	novelRoute.Post("/create", createNovel(db))
	novelRoute.Post("/from/:username", getUsersNovels(db))
	novelRoute.Post("/:novelID", getNovel(db))

	novelRoute.Patch("/:novelID", updateNovelMetadata(db))

	novelRoute.Delete("/:novelID", deleteNovel(db))
}

// Get Novel
//
//	@Summary		Get the novel with provided novel id
//	@Description	If the novel is private, the user need to be logged in with the author account
//	@Tags			novel
//	@Produce		json
//	@Param			NovelID			path		string						true	"Novel ID"
//	@Param			sessionString	body		model.IncludeSessionString	false	"User's Session"
//	@Success		200				{object}	model.NovelView
//	@Failure		404
//	@Failure		500
//	@Router			/novel/:novelID [POST]
func getNovel(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		novelIDStr := c.Params("novelID")
		if len(novelIDStr) != model.IDHexLength {
			return c.SendStatus(fiber.StatusNotFound)
		}
		novelID, err := Unhex(novelIDStr)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		novelView, ok := db.GetNovelView(novelID)
		if !ok {
			return c.SendStatus(fiber.StatusNotFound)
		}

		if novelView.Visibility == model.VisibilityPrivate.String() {
			if c.Locals(middleware.KeyIsUserAuth) == false {
				return c.SendStatus(fiber.StatusUnauthorized)
			}
			session, ok := c.Locals(middleware.KeyUserSession).(model.Session)
			if !ok {
				log.Warn("Check auth middleware")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			if !ok {
				return c.SendStatus(fiber.StatusUnauthorized)
			}

			if hex.EncodeToString(session.UserID) != novelView.Author.ID {
				return c.SendStatus(fiber.StatusUnauthorized)
			}
		}

		// Everything is good
		return c.JSON(novelView)
	}
}

type createNovelResult struct {
	NovelID string `json:"novel_id"`
}

// Create Novel
//
//	@Summary		Create a new novel with the provided metadata, return the created novel id
//	@Description	Possible error code: MissingField, InvalidLanguageFormat, TitleTooLong, TaglineTooLong
//	@Tags			novel
//	@Accept			json
//	@Produce		json
//	@Param			NovelDetails	body		model.NovelMetadata			true	"Novel details"
//	@Param			sessionString	body		model.IncludeSessionString	true	"User's Session"
//	@Success		201				{object}	createNovelResult
//	@Failure		400				{object}	ErrorJSON
//	@Failure		401
//	@Failure		500
//	@Router			/novel/create [POST]
func createNovel(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals(middleware.KeyIsUserAuth) == false {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		var input model.NovelMetadata
		err := c.BodyParser(&input)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(BadInput))
		}

		if ok, code := checkNovelMetadata(&input); !ok {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(code))
		}

		session, ok := c.Locals(middleware.KeyUserSession).(model.Session)
		if !ok {
			log.Warn("Check the authentication middleware")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		input.Author = session.UserID
		uid, ok := db.CreateNovel(&input)
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusCreated).JSON(
			createNovelResult{
				NovelID: hex.EncodeToString(uid),
			})
	}
}

// Update Novel Metadata
//
//	@Summary		Update the novel metadata with the provided metadata
//	@Description	Possible error code: MissingField, InvalidLanguageFormat, TitleTooLong, TaglineTooLong
//	@Tags			novel
//	@Accept			json
//	@Param			NovelID			path	string						true	"Novel ID"
//	@Param			NovelDetails	body	model.NovelMetadata			true	"Novel details"
//	@Param			sessionString	body	model.IncludeSessionString	true	"User's Session"
//	@Success		200
//	@Failure		400	{object}	ErrorJSON
//	@Failure		401
//	@Failure		404
//	@Failure		500
//	@Router			/novel/:novelID [PATCH]
func updateNovelMetadata(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals(middleware.KeyIsUserAuth) == false {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		var input model.NovelMetadata
		err := c.BodyParser(&input)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(BadInput))
		}

		if ok, code := checkNovelMetadata(&input); !ok {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(code))
		}

		novelIDStr := c.Params("novelID")
		if len(novelIDStr) != model.IDHexLength {
			return c.SendStatus(fiber.StatusNotFound)
		}
		novelID, err := Unhex(novelIDStr)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		session, ok := c.Locals(middleware.KeyUserSession).(model.Session)
		if !ok {
			log.Warn("Check the authentication middleware")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if bytes.Compare(session.UserID, input.Author) != 0 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		ok = db.UpdateNovelMetadata(novelID, &input)
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

// Get User's Novels
//
//	@Summary		Get all the novels from the user with the provided user id
//	@Description	If the user is not logged in, only the public novels will be returned
//	@Tags			novel
//	@Produce		json
//	@Param			username		path		string						true	"username"
//	@Param			sessionString	body		model.IncludeSessionString	true	"User's Session"
//	@Param			filtersAndSort	query		model.FiltersAndSortNovel	false	"Filters and sorting options"
//	@Success		200				{object}	[]model.NovelMetadataSmall
//	@Failure		401
//	@Failure		404
//	@Failure		500
//	@Router			/novel/from/:username [POST]
func getUsersNovels(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Params("username")
		if !IsUsernameValid(username) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		user, ok := db.GetUser(username)
		if !ok {
			return c.SendStatus(fiber.StatusNotFound)
		}
		filtersAndSort := getFiltersAndSort(c)

		isSelf := c.Locals(middleware.KeyIsUserAuth) == true
		if isSelf {
			session, _ := c.Locals(middleware.KeyUserSession).(model.Session)
			isSelf = bytes.Compare(session.UserID, user.ID) == 0
		}

		novelsMetadataSmall := db.GetUsersNovels(user.ID, &filtersAndSort, isSelf)
		return c.JSON(novelsMetadataSmall)
	}
}

// Search and Filter Novels
//
//	@Summary	Search and filter novels with the provided filters and sorting options, if no filters and sorting options are provided, all the public novels will be returned
//	@Tags		novel
//	@Produce	json
//	@Param		filtersAndSort	query		model.FiltersAndSortNovel	false	"Filters and sorting options"
//	@Success	200				{object}	[]model.NovelMetadataSmall
//	@Failure	500
//	@Router		/novel/find [GET]
func searchAndFilterNovel(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filtersAndSortOption := getFiltersAndSort(c)

		novels := db.FindNovels(&filtersAndSortOption)
		return c.JSON(novels)
	}
}

// Delete Novel
//
//	@Deprecated
//	@Summary	Delete the novel and all the related stuff like volumes, chapters, comments, images with the provided novel id
//	@Tags		novel
//	@Param		NovelID	path	string	true	"Novel ID"
//	@Success	200
//	@Failure	401
//	@Failure	404
//	@Failure	500
//	@Router		/novel/:novelID [DELETE]
func deleteNovel(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
		//novelID := c.Params("novelID")
		//if len(novelID) != model.IDHexLength {
		//	return c.SendStatus(fiber.StatusNotFound)
		//}
		//
		//if c.Locals(middleware.KeyIsUserAuth) == false {
		//	return c.SendStatus(fiber.StatusUnauthorized)
		//}
		//
		//session, ok := c.Locals(middleware.KeyUserSession).(model.Session)
		//if !ok {
		//	log.Warn("Check the authentication middleware")
		//	return c.SendStatus(fiber.StatusInternalServerError)
		//}
		//
		//novel, ok := db.GetNovel(novelID)
		//if !ok {
		//	return c.SendStatus(fiber.StatusNotFound)
		//}
		//
		//if bytes.Compare(session.UserID, novel.Author) != 0 {
		//	return c.SendStatus(fiber.StatusUnauthorized)
		//}
		//
		//ok = db.DeleteNovel(novelID)
		//if !ok {
		//	return c.SendStatus(fiber.StatusInternalServerError)
		//}
		//
		//return c.SendStatus(fiber.StatusOK)
	}
}

func checkNovelMetadata(input *model.NovelMetadata) (bool, ErrorCode) {
	if matched, err := regexp.Match("^[a-z]{3}$", []byte(input.Language)); err != nil ||
		matched == false {
		return false, InvalidLanguageFormat
	}

	if utf8.RuneCountInString(input.Title) > model.TitleMaxLength {
		return false, TitleTooLong
	}

	if utf8.RuneCountInString(input.Tagline) > model.TaglineMaxLength {
		return false, TaglineTooLong
	}

	if utf8.RuneCountInString(input.Description) > model.DescriptionMaxLength {
		return false, DescriptionTooLong
	}

	if input.Visibility.String() == "Unknown" || input.Status.String() == "Unknown" {
		return false, BadInput
	}

	return true, BadInput
}
