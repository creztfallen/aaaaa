package main

import (
	"github.com/creztfallen/compiler_api/cmd/api/config"
	"github.com/creztfallen/compiler_api/cmd/api/users"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New()
	config.ConnectDB()

	users.UserRoute(app)

	err := app.Listen(":6000")
	if err != nil {
		log.Fatal(err)
	}
}
