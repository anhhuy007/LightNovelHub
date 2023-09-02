package main

import (
	"Lightnovel/middleware"
	"Lightnovel/model"
	"Lightnovel/route"
	"context"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

//	@title		Light novel API
//	@version	1.0

//	@host		localhost:8080
//	@BasePath	/api/v1
func main() {
	mysqlConfig := mysql.Config{
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Addr:      os.Getenv("MYSQL_HOST"),
		Net:       "tcp",
		DBName:    os.Getenv("MYSQL_DATABASE"),
		ParseTime: true,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	log.Debug(mysqlConfig.FormatDSN())
	db, err := sqlx.ConnectContext(ctx, "mysql", mysqlConfig.FormatDSN())
	cancel()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Error(err)
		}
	}()
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

	swagger := swagger.New(swagger.Config{
		FilePath: "./docs/swagger.json",
		BasePath: "/",
	})
	app.Use(swagger)

	app.All("/ok", func(c *fiber.Ctx) error {
		log.Debug(c.Locals(middleware.KeyIsUserAuth))
		return c.SendString("Hello, World 👋!")
	})

	apiRoute := app.Group("/api")
	v1 := apiRoute.Group("/v1")

	route.AddAccountRoutes(&v1, &database)
	route.AddUploadRoutes(&v1, &database)

	log.Fatal(app.Listen(":8080"))
}
