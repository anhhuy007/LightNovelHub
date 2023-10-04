package main

import (
	"Lightnovel/middleware"
	"Lightnovel/model/repo"
	"Lightnovel/route"
	"context"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

//	@title		Light novel API
//	@version	1.0

// @BasePath	/api/v1
func main() {
	db := getDBConnect()
	defer func() {
		err := db.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	database := repo.NewDatabase(db, time.Minute)

	app := fiber.New()

	//file, err := os.Create(fmt.Sprintf("logs/%v.txt", time.Now().Format("2006-01-02-15-04-05")))
	//if err != nil {
	//	panic(err)
	//}
	//defer func() {
	//	err := file.Close()
	//	if err != nil {
	//		log.Error(err)
	//	}
	//}()
	//app.Use(logger.New(logger.Config{
	//	Output:   file,
	//	TimeZone: "UTC",
	//}))

	app.Use(logger.New(logger.ConfigDefault))
	app.Get("/metrics", monitor.New())
	authMiddleware := middleware.AuthenticationCheck(&database)
	app.Use(authMiddleware)

	swaggerRoute := app.Group("/swagger")
	swaggerRoute.Use(swagger.New(swagger.Config{
		BasePath: "/swagger",
		FilePath: "./docs/swagger.json",
	}))

	app.Get("/ok", func(c *fiber.Ctx) error {
		log.Debug(c.Locals(middleware.KeyIsUserAuth))
		return c.SendString("Hello, World 👋!")
	})

	apiRoute := app.Group("/api")
	v1 := apiRoute.Group("/v1")

	route.AddAccountRoutes(&v1, &database)
	route.AddUploadRoutes(&v1, &database)

	//data, _ := json.MarshalIndent(app.Stack(), "", "  ")
	//fmt.Println(string(data))

	log.Fatal(app.Listen(":8080"))
}

func getDBConnect() *sqlx.DB {
	mysqlConfig := mysql.Config{
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Addr:      os.Getenv("MYSQL_HOST"),
		Net:       "tcp",
		DBName:    os.Getenv("MYSQL_DATABASE"),
		ParseTime: true,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	//log.Debug(mysqlConfig.FormatDSN())
	db, err := sqlx.ConnectContext(ctx, "mysql", mysqlConfig.FormatDSN())
	cancel()
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if db.Ping() != nil {
		panic(err)
	}
	return db
}
