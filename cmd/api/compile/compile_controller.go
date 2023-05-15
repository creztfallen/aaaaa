package compile

import (
	"github.com/creztfallen/compiler_api/cmd/api/responses"
	"github.com/creztfallen/compiler_api/cmd/api/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	counter int
	mutex   sync.Mutex
)

//var userCollection = config.GetCollection(config.DB, "users")

func Compile(c *fiber.Ctx) error {

	mutex.Lock()
	counter++
	mutex.Unlock()

	err := utils.UploadFileHandler(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DefaultResponse{
			Status:  http.StatusInternalServerError,
			Message: "couldn't save the file",
			Data:    nil})
	}

	time.Sleep(5 * time.Second)

	filePath := "./cmd/api/compile/release/P24_HGUV6.5.3B01_UPGRADE.bin"
	content, err := os.ReadFile(filePath)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DefaultResponse{
			Status:  http.StatusInternalServerError,
			Message: "Couldn't read the file",
			Data:    nil})
	}

	_ = utils.S3Upload(c, content, counter)

	objectURL, err := utils.S3Get(counter)

	return c.Status(http.StatusCreated).JSON(responses.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &fiber.Map{"firmware": objectURL}})
}
