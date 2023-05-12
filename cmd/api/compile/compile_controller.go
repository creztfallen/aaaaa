package compile

import (
	"github.com/creztfallen/compiler_api/cmd/api/responses"
	"github.com/creztfallen/compiler_api/cmd/api/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func Compile(c *fiber.Ctx) error {
	err := utils.UploadFileHandler(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DefaultResponse{
			Status:  http.StatusInternalServerError,
			Message: "file saved",
			Data:    nil})
	}
	return c.Status(http.StatusCreated).JSON(responses.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "file saved",
		Data:    nil})
}
