package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type AuthCredentials struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	DeviceName string `json:"deviceName"`
}

type AuthReturn struct {
	SessionID string `json:"sessionID"`
}

func AddAuthRoutes(app *fiber.App, db *Database) {
	app.Post("/auth/login", func(c *fiber.Ctx) error {
		var authCredentials AuthCredentials
		err := c.BodyParser(&authCredentials)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		user, err := db.GetUser(authCredentials.Username)
		if err != nil {
			// Do the hash comparison anyway to prevent timing attacks
			passwordVerify(authCredentials.Password, user.Password[:])
			return c.SendStatus(fiber.StatusNotFound)
		}
		if !passwordVerify(authCredentials.Password, user.Password[:]) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		sessionID, err := db.CreateSession(user.ID, authCredentials.DeviceName)
		if err != nil {
			log.Debug(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(AuthReturn{sessionID})
	})

	app.Post("/auth/register", func(c *fiber.Ctx) error {
		var authCredentials AuthCredentials
		err := c.BodyParser(&authCredentials)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		userId, err := db.CreateUser(authCredentials.Username, authCredentials.Password)
		if err != nil {
			log.Debug(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		sessionId, err := db.CreateSession(userId, authCredentials.DeviceName)
		if err != nil {
			log.Warn(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(AuthReturn{sessionId})
	})
}
