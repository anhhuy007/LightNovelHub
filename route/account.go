package route

import (
	"Lightnovel/model"
	"Lightnovel/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type authCredentials struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	DeviceName string `json:"deviceName"`
}

func AddAccountRoutes(router *fiber.Router, db model.DB) {
	accountRoute := (*router).Group("/account")

	accountRoute.Post("/login", func(c *fiber.Ctx) error {
		var authCredentials authCredentials
		err := c.BodyParser(&authCredentials)
		log.Debugf("%#v", authCredentials)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		user, err := db.GetUser(authCredentials.Username)
		// Do the hash comparison anyway to prevent timing attacks
		if errors.Is(err, model.ErrUserNotFound) {
			utils.PasswordVerify(authCredentials.Password, []byte{})
			return c.SendStatus(fiber.StatusNotFound)
		} else if err != nil {
			log.Error(err)
			utils.PasswordVerify(authCredentials.Password, []byte{})
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if !utils.PasswordVerify(authCredentials.Password, user.Password) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		sessionInfo, err := db.CreateSession(
			user.ID,
			authCredentials.DeviceName,
		)

		if err != nil {
			log.Warn(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(sessionInfo)
	})

	accountRoute.Post("/register", func(c *fiber.Ctx) error {
		var authCredentials authCredentials
		err := c.BodyParser(&authCredentials)
		log.Debugf("%#v", authCredentials)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		userId, err := db.CreateUser(
			authCredentials.Username,
			authCredentials.Password,
		)
		if errors.Is(err, model.ErrInvalidPassword) || errors.Is(err, model.ErrUserAlreadyExist) {
			return c.SendStatus(fiber.StatusBadRequest)
		} else if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		sessionInfo, err := db.CreateSession(userId, authCredentials.DeviceName)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(sessionInfo)
	})

	accountRoute.Post("/logout", func(c *fiber.Ctx) error {
		var sessionStr model.IncludeSessionString
		err := c.BodyParser(&sessionStr)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		_ = db.DeleteSession(sessionStr.Session)
		return c.SendStatus(fiber.StatusOK)
	})

	accountRoute.Post("/renew", func(c *fiber.Ctx) error {
		var oldSessionStr model.IncludeSessionString
		err := c.BodyParser(&oldSessionStr)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		session, err := db.GetSession(oldSessionStr.Session)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		newSessionInfo, err := db.CreateSession(session.UserID, session.DeviceName)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		_ = db.DeleteSession(oldSessionStr.Session)
		return c.JSON(newSessionInfo)
	})
}
