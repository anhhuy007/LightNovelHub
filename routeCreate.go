package main

import "github.com/gofiber/fiber/v2"

func AddUploadRoutes(app *fiber.App, db *Database) {
	app.Post("/create/novel")
}
