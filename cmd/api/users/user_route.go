package users

import "github.com/gofiber/fiber/v2"

func UserRoute(app *fiber.App) {
	app.Post("/signup", CreateUser)
	app.Post("/signin", Login)
}
