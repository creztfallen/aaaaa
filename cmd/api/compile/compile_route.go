package compile

import "github.com/gofiber/fiber/v2"

func CompileRoute(app *fiber.App) {
	app.Post("/compile", Compile)
}
