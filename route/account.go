package route

import (
	"Lightnovel/middleware"
	"Lightnovel/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"unicode/utf8"
)

// TODO: Delete User
func AddAccountRoutes(router *fiber.Router, db model.DB) {
	accountRoute := (*router).Group("/accounts")

	accountRoute.Post("/login", login(db))

	accountRoute.Post("/register", register(db))

	accountRoute.Post("/logout", logout(db))

	accountRoute.Post("/renew", renew(db))

	accountRoute.Delete("/:username", deleteUser(db))

	accountRoute.Patch("/update", updateUser(db))

	accountRoute.Get("/:username", getUserView(db))

	accountRoute.Post("/changepassword", changeUserPassword(db))

	accountRoute.Post("/self", getUserViewFromSession(db))

	accountRoute.Post("/followed/users", getFollowedUser(db))

	accountRoute.Post("/followed/novels", getFollowedNovel(db))
}

// Login
//
//	@Summary		Log the user in, return a new user session
//	@Description	The session token should be renewed a week before expires, possible error: WrongPassword, UserNotFound, BadInput, BadPassword, BadUsername, BadDeviceName
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			credential	body		authCredentials	true	"User credentials"
//	@Success		200			{object}	model.SessionInfo
//	@Failure		400			{object}	ErrorJSON
//	@Failure		500
//	@Router			/accounts/login [POST]
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

		hashedPassword, err := PasswordHash(authCredentials.Password)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(BadPassword))
		}
		userId, ok := db.CreateUser(
			authCredentials.Username,
			hashedPassword,
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
		session, err := Unhex(sessionStr.Session)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		_ = db.DeleteSession(session)
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

		oldSession, err := Unhex(oldSessionStr.Session)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		session, ok := db.GetSession(oldSession)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		newSessionInfo, ok := db.CreateSession(session.UserID, session.DeviceName)
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		_ = db.DeleteSession(oldSession)
		return c.JSON(newSessionInfo)
	}
}

// Delete User
//
//	@Deprecated
//	@Summary		Delete user's account and all other data
//	@Description	Possible error: BadInput, BadPassword, BadUsername, UserNotFound
//	@Tags			accounts
//	@Accept			json
//	@Param			credential	body	requiredCredential	true	"User credentials"
//	@Success		200
//	@Failure		400	{object}	ErrorJSON
//	@Failure		500
//	@Router			/accounts/delete [DELETE]
func deleteUser(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
		//var credential requiredCredential
		//err := c.BodyParser(&credential)
		//if err != nil {
		//	return c.SendStatus(fiber.StatusBadRequest)
		//}
		//
		//if ok, code := credential.Validate(); !ok {
		//	return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(code))
		//}
		//
		//user, ok := db.GetUser(credential.Username)
		//if !ok {
		//	return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(UserNotFound))
		//}
		//
		//passwordGood := PasswordVerify(credential.Password, user.Password)
		//if !passwordGood {
		//	return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(WrongPassword))
		//}
		//
		//db.DeletaAllSessions(user.ID)
		//ok = db.DeleteUser(user.ID)
		//if !ok {
		//	return c.SendStatus(fiber.StatusInternalServerError)
		//}
		//return c.SendStatus(fiber.StatusOK)
	}
}

// Update User
//
//	@Summary		Update user's metadata
//	@Description	Possible error: BadInput, BadUsername, BadDisplayname, BadEmail, UserAlreadyExists
//	@Tags			accounts
//	@Accept			json
//	@Param			credential	body	model.UserMetadata	true	"User metadata"
//	@Success		200
//	@Failure		400	{object}	ErrorJSON
//	@Failure		401
//	@Failure		500
//	@Router			/accounts/update [PATCH]
func updateUser(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals(middleware.KeyIsUserAuth) == false {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		var input model.UserMetadata
		err := c.BodyParser(&input)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(BadInput))
		}

		if ok, code := checkUserMetadata(input); !ok {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(code))
		}

		session, ok := c.Locals(middleware.KeyUserSession).(model.Session)
		if !ok {
			log.Warn("Check the authentication middleware")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		_, ok = db.GetUser(input.Username)
		if ok {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(UserAlreadyExists))
		}
		ok = db.UpdateUserMetadata(session.UserID, input)
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

// Get User
//
//	@Summary	Get user's metadata
//	@Tags		accounts
//	@Param		userID	path		string	true	"UserId"
//	@Success	200		{object}	model.UserView
//	@Failure	404
//	@Failure	500
//	@Router		/accounts/:username [GET]
func getUserView(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Params("username")

		userView, ok := db.GetUserView(username)
		if !ok {
			return c.SendStatus(fiber.StatusNotFound)
		}

		return c.JSON(userView)
	}
}

type changePasswordCredential struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (cr *changePasswordCredential) Validate() (bool, ErrorCode) {
	if !(IsPasswordValid(cr.OldPassword) && IsPasswordValid(cr.NewPassword)) {
		return false, BadPassword
	}

	return true, BadInput
}

// Change Password
//
//	@Summary		Change user's password
//	@Description	Possible error: BadInput, BadPassword, WrongPassword
//	@Tags			accounts
//	@Accept			json
//	@Param			credential	body	changePasswordCredential	true	"Old and new password"
//	@Success		200
//	@Failure		400	{object}	ErrorJSON
//	@Failure		401
//	@Failure		500
//	@Router			/accounts/changepassword [POST]
func changeUserPassword(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals(middleware.KeyIsUserAuth) == false {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		var input changePasswordCredential
		err := c.BodyParser(&input)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(BadInput))
		}

		if ok, code := input.Validate(); !ok {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(code))
		}

		session, ok := c.Locals(middleware.KeyUserSession).(model.Session)
		if !ok {
			log.Warn("Check the authentication middleware")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		user, ok := db.GetUserWithID(session.UserID)
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		newHashed, err := PasswordHash(input.NewPassword)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(BadPassword))
		}

		if PasswordVerify(input.OldPassword, user.Password) == false {
			return c.Status(fiber.StatusBadRequest).JSON(buildErrorJSON(WrongPassword))
		}
		db.UpdateUserPassword(user.ID, newHashed)
		return c.SendStatus(fiber.StatusOK)
	}
}

// Get Self
//
//	@Summary	Get user's metadata from session
//	@Tags		accounts
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	model.UserView
//	@Failure	401
//	@Failure	500
//	@Router		/accounts/self [POST]
func getUserViewFromSession(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals(middleware.KeyIsUserAuth) == false {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		session, ok := c.Locals(middleware.KeyUserSession).(model.Session)
		if !ok {
			log.Warn("Check the authentication middleware")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		UserView, ok := db.GetUserViewWithID(session.UserID)
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(UserView)
	}
}

// Get Followed Users
//
//	@Summary	Get user's followed users
//	@Tags		accounts
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]model.UserMetadataSmall
//	@Failure	401
//	@Failure	500
//	@Router		/accounts/followed/users [POST]
func getFollowedUser(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals(middleware.KeyIsUserAuth) == false {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		session, ok := c.Locals(middleware.KeyUserSession).(model.Session)
		if !ok {
			log.Warn("Check the authentication middleware")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(db.GetFollowedUser(session.UserID))
	}
}

// Get Followed Novels
//
//	@Summary	Get user's followed novels
//	@Tags		accounts
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]model.NovelMetadataSmall
//	@Failure	401
//	@Failure	500
//	@Router		/accounts/followed/novels [POST]
func getFollowedNovel(db model.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals(middleware.KeyIsUserAuth) == false {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		session, ok := c.Locals(middleware.KeyUserSession).(model.Session)
		if !ok {
			log.Warn("Check the authentication middleware")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(db.GetFollowedNovel(session.UserID))
	}
}

type requiredCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (rc *requiredCredential) Validate() (bool, ErrorCode) {
	if IsUsernameValid(rc.Username) == false {
		return false, BadUsername
	}

	if IsPasswordValid(rc.Password) == false {
		return false, BadPassword
	}

	return true, BadInput
}

type authCredentials struct {
	requiredCredential
	DeviceName string `json:"deviceName"`
}

func (ac *authCredentials) Validate() (bool, ErrorCode) {
	if utf8.RuneCountInString(ac.DeviceName) > model.DeviceNameMaxLength {
		return false, BadDeviceName
	}

	return ac.requiredCredential.Validate()
}
