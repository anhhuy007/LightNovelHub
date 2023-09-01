package middleware

import (
	"Lightnovel/model"
	"github.com/gofiber/fiber/v2"
	"time"
)

const (
	KeyIsUserAuth  = "isUserAuth"
	KeyUserSession = "userSession"
)

func AuthenticationCheck(db model.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		body := map[string]string{}
		err := c.BodyParser(&body)
		if err != nil {
			c.Locals(KeyIsUserAuth, false)
			return c.Next()
		}

		sessionInfoStr, ok := body["session"]
		if !ok {
			c.Locals(KeyIsUserAuth, false)
			return c.Next()
		}
		session, err := db.GetSession(sessionInfoStr)
		if err != nil || session.ExpireAt.Before(time.Now()) {
			c.Locals(KeyIsUserAuth, false)
			return c.Next()
		}

		c.Locals(KeyIsUserAuth, true)
		c.Locals(KeyUserSession, session)
		return c.Next()
	}

}
