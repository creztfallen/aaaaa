package compile

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/creztfallen/compiler_api/cmd/api/config"
	"github.com/creztfallen/compiler_api/cmd/api/responses"
	"github.com/creztfallen/compiler_api/cmd/api/utils"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("sa-east-1"),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(config.AWS_KEY(), config.AWS_SECRET_KEY(), ""),
	})
	if err != nil {
		return err
	}

	s3Client := s3.New(sess)

	err = utils.UploadFileHandler(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DefaultResponse{
			Status:  http.StatusInternalServerError,
			Message: "couldn't save the file",
			Data:    nil})
	}

	output, _ := utils.ExecChain("sentinel")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DefaultResponse{
			Status:  http.StatusInternalServerError,
			Message: "Unexpected error while dealing with the sentinel",
			Data:    &fiber.Map{"error": output}})
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

	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(config.AWS_BUCKET()),
		Key:    aws.String(strconv.Itoa(counter) + "-P24_HGUV6.5.3B01_UPGRADE.bin"),
		Body:   bytes.NewReader(content),
	}

	_, err = s3Client.PutObject(uploadInput)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DefaultResponse{
			Status:  http.StatusInternalServerError,
			Message: "Couldn't send the file to bucket",
			Data:    &fiber.Map{"error": err}})
	}

	params := &s3.GetObjectInput{
		Bucket: aws.String(config.AWS_BUCKET()),
		Key:    aws.String(strconv.Itoa(counter) + "-P24_HGUV6.5.3B01_UPGRADE.bin"),
	}

	req, _ := s3Client.GetObjectRequest(params)
	objectURL, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Fatal("[AWS GET LINK]:", params, err)
	}

	return c.Status(http.StatusCreated).JSON(responses.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &fiber.Map{"firmware": objectURL}})
}
