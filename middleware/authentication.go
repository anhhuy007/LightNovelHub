package middleware

import (
	"Lightnovel/model"
	"encoding/hex"
	"github.com/gofiber/fiber/v2"
	"time"
)

const (
	KeyIsUserAuth  = "isUserAuth"
	KeyUserSession = "userSession"
	BodySession    = "session"
)

func Unhex(s string) ([]byte, error) {
	return hex.DecodeString(s[:model.IDHexLength])
}

func AddAuthenticationCheck(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body model.IncludeSessionString
		err := c.BodyParser(&body)
		if err != nil || len(body.Session) != model.IDHexLength {
			c.Locals(KeyIsUserAuth, false)
			return c.Next()
		}

		sessionInfo, err := Unhex(body.Session)
		if err != nil {
			c.Locals(KeyIsUserAuth, false)
			return c.Next()
		}
		session, ok := db.GetSession(sessionInfo)
		if !ok || session.ExpireAt.Before(time.Now()) {
			c.Locals(KeyIsUserAuth, false)
			return c.Next()
		}

		c.Locals(KeyIsUserAuth, true)
		c.Locals(KeyUserSession, session)
		return c.Next()
	}
}
