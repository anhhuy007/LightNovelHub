package main

import (
	"Lightnovel/middleware"
	"Lightnovel/model"
	"Lightnovel/route"
	"context"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	mysqlConfig := mysql.Config{
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Addr:      os.Getenv("MYSQL_HOST"),
		Net:       "tcp",
		DBName:    os.Getenv("MYSQL_DATABASE"),
		ParseTime: true,
	}
	log.Debug(mysqlConfig.FormatDSN())
	db, err := sqlx.ConnectContext(ctx, "mysql", mysqlConfig.FormatDSN())
	cancel()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if db.Ping() != nil {
		panic(err)
	}
	database := model.NewDatabase(db, time.Minute)

	app := fiber.New()
	authMiddleware := middleware.AuthenticationCheck(&database)
	app.Use(authMiddleware)

	app.All("/", func(c *fiber.Ctx) error {
		log.Debug(c.Locals(middleware.KeyIsUserAuth))
		return c.SendString("Hello, World 👋!")
	})

	apiRoute := app.Group("/api")
	v1 := apiRoute.Group("/v1")

	route.AddAccountRoutes(&v1, &database)
	route.AddUploadRoutes(&v1, &database)

	log.Fatal(app.Listen(":8080"))
}
