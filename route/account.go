package route

import (
	"Lightnovel/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"unicode/utf8"
)

type authCredentials struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	DeviceName string `json:"deviceName"`
}

func (ac *authCredentials) Validate() (bool, ErrorCode) {
	if IsUsernameValid(ac.Username) == false {
		return false, BadUsername
	}

	if IsPasswordValid(ac.Password) == false {
		return false, BadPassword
	}

	if utf8.RuneCountInString(ac.DeviceName) > model.DeviceNameMaxLength {
		return false, BadDeviceName
	}

	return true, BadInput
}

func AddAccountRoutes(router *fiber.Router, db model.DB) {
	accountRoute := (*router).Group("/accounts")

	accountRoute.Post("/login", login(db))

	accountRoute.Post("/register", register(db))

	accountRoute.Post("/logout", logout(db))

	accountRoute.Post("/renew", renew(db))
}

//	Login
//
// @Summary		Log the user in, return a new user session
// @Description	The session token should be renewed a week before expires, possible error: WrongPassword, UserNotFound, BadInput, BadPassword, BadUsername, BadDeviceName
// @Tags			accounts
// @Accept			json
// @Produce		json
// @Param			credential	body		authCredentials	true	"User credentials"
// @Success		200			{object}	model.SessionInfo
// @Failure		400			{object}	ErrorJSON
// @Failure		500
// @Router			/accounts/login [POST]
func login(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var authCredentials authCredentials
		err := c.BodyParser(&authCredentials)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(BadInput))
		}

		if ok, code := authCredentials.Validate(); !ok {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(code))
		}

		user, ok := db.GetUser(authCredentials.Username)
		// Do the hash comparison anyway to prevent timing attacks
		passwordGood := PasswordVerify(authCredentials.Password, user.Password)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(UserNotFound))
		}

		if !passwordGood {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(WrongPassword))
		}

		sessionInfo, ok := db.CreateSession(
			user.ID,
			authCredentials.DeviceName,
		)

		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(sessionInfo)
	}
}

// Register
//
//	@Summary		Register the user, return a new user session
//	@Description	Possible error: BadInput, BadPassword, BadUsername, BadDeviceName, UserAlreadyExists
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			credential	body		authCredentials	true	"User credentials"
//	@Success		201			{object}	model.SessionInfo
//	@Failure		400			{object}	ErrorJSON
//	@Failure		500
//	@Router			/accounts/register [POST]
func register(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var authCredentials authCredentials
		err := c.BodyParser(&authCredentials)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if ok, code := authCredentials.Validate(); !ok {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(code))
		}

		_, ok := db.GetUser(authCredentials.Username)
		if ok {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(UserAlreadyExists))
		}
		userId, ok := db.CreateUser(
			authCredentials.Username,
			authCredentials.Password,
		)
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		sessionInfo, ok := db.CreateSession(userId, authCredentials.DeviceName)
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.Status(fiber.StatusCreated).JSON(sessionInfo)
	}
}

// Logout
//
//	@Summary	Log the user out
//	@Tags		accounts
//	@Accept		json
//	@Param		credential	body	model.IncludeSessionString	true	"User credentials"
//	@Success	200
//	@Failure	400
//	@Router		/accounts/logout [POST]
func logout(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var sessionStr model.IncludeSessionString
		err := c.BodyParser(&sessionStr)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		_ = db.DeleteSession(sessionStr.Session)
		return c.SendStatus(fiber.StatusOK)
	}
}

// Renew
//
//	@Summary	Renew the session token, the token should be renewed a week before expires
//	@Tags		accounts
//	@Accept		json
//	@Produce	json
//	@Param		credential	body		model.IncludeSessionString	true	"User credentials"
//	@Success	200			{object}	model.SessionInfo
//	@Failure	400
//	@Failure	401
//	@Failure	500
//	@Router		/accounts/renew [POST]
func renew(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var oldSessionStr model.IncludeSessionString
		err := c.BodyParser(&oldSessionStr)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		session, ok := db.GetSession(oldSessionStr.Session)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		newSessionInfo, ok := db.CreateSession(session.UserID, session.DeviceName)
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		_ = db.DeleteSession(oldSessionStr.Session)
		return c.JSON(newSessionInfo)
	}
}

func PasswordVerify(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}

func IsUsernameValid(username string) bool {
	matched, _ := regexp.MatchString(fmt.Sprintf("^[a-zA-Z]{%v,%v}$", model.UserNameMinLength, model.UserNameMaxLength), username)
	return matched
}

func IsPasswordValid(password string) bool {
	return !(len(password) < model.PasswordMinLength || len(password) > model.PasswordMaxLength)
}
