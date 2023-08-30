package main

import (
	"context"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	mysqlConfig := mysql.Config{
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASSWORD"),
		Addr:   os.Getenv("MYSQL_HOST"),
		Net:    "tcp",
		DBName: os.Getenv("MYSQL_DATABASE"),
	}

	db, err := sqlx.ConnectContext(ctx, "mysql", mysqlConfig.FormatDSN())
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

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Listen(":8080")
}
