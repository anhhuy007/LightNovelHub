package middleware

import (
	"Lightnovel/model"
	"github.com/gofiber/fiber/v2"
	"time"
)

const (
	KeyIsUserAuth  = "isUserAuth"
	KeyUserSession = "userSession"
	BodySession    = "session"
)

func AuthenticationCheck(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := map[string]string{}
		err := c.BodyParser(&body)
		if err != nil {
			c.Locals(KeyIsUserAuth, false)
			return c.Next()
		}

		sessionInfoStr, ok := body[BodySession]
		if !ok {
			c.Locals(KeyIsUserAuth, false)
			return c.Next()
		}
		session, ok := db.GetSession(sessionInfoStr)
		if !ok || session.ExpireAt.Before(time.Now()) {
			c.Locals(KeyIsUserAuth, false)
			return c.Next()
		}

		c.Locals(KeyIsUserAuth, true)
		c.Locals(KeyUserSession, session)
		return c.Next()
	}

}
