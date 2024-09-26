package main

import (
    "log"
	"flag"
	"strconv"
	"hudaisoft-backend/lib"
    "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger")

var (
	DatabasePath string 
	WorkDirectory string
	Port int
)

func main() {
	// Get parameters.
	flag.StringVar(&DatabasePath, "databasepath", "hudautomata.sqlite", "Database location, using SQLite3 for this app.")
	flag.StringVar(&WorkDirectory, "workdirectory", "hudautomata", "Asset's and all other thing will goes that directory, default is: 'hudautomata'.")
	flag.IntVar(&Port, "port", 3000, "Application port which will run's in there, default is: '3000'.")

	// Setup and init.		
	flag.Parse()
	app := fiber.New()

	app.Use(logger.New())

	lib.SetupDatabase(DatabasePath)

	// Backend Api.
    app.Get("/", func (c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

	// For begining we don't need to add custom interface.
	port := "0.0.0.0:" + strconv.Itoa(Port)
    log.Fatal(app.Listen(port))
}
