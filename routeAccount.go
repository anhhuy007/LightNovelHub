package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type authCredentials struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	DeviceName string `json:"deviceName"`
}

func AddAuthRoutes(app *fiber.App, db *Database) {

	app.Post("/auth/login", func(c *fiber.Ctx) error {
		var authCredentials authCredentials
		err := c.BodyParser(&authCredentials)
		log.Debugf("%#v", authCredentials)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		user, err := db.GetUser(authCredentials.Username)
		// Do the hash comparison anyway to prevent timing attacks
		if err == ErrUserNotFound {
			passwordVerify(authCredentials.Password, []byte{})
			return c.Status(fiber.StatusNotFound).JSON(ErrorMessage{
				err.Error(),
			})
		} else if err != nil {
			log.Error(err)
			passwordVerify(authCredentials.Password, []byte{})
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if !passwordVerify(authCredentials.Password, user.Password) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorMessage{
				ErrWrongPassword.Error(),
			})
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

	app.Post("/auth/register", func(c *fiber.Ctx) error {
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
		if err == ErrInvalidPassword || err == ErrUserAlreadyExist {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorMessage{
				err.Error(),
			})
		} else if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		sessionInfo, err := db.CreateSession(userId, authCredentials.DeviceName)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(sessionInfo)
	})

	app.Post("/auth/logout", func(c *fiber.Ctx) error {
		var sessionStr IncludeSessionString
		err := c.BodyParser(&sessionStr)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		_ = db.DeleteSession(sessionStr.Session)
		return c.SendStatus(fiber.StatusOK)
	})

	app.Post("/auth/renew", func(c *fiber.Ctx) error {
		var oldSessionStr IncludeSessionString
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
