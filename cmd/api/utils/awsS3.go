package utils

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/creztfallen/compiler_api/cmd/api/config"
	"github.com/creztfallen/compiler_api/cmd/api/responses"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"strconv"
	"time"
)

func S3Upload(c *fiber.Ctx, content []byte, counter int) error {

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("sa-east-1"),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(config.AWS_KEY(), config.AWS_SECRET_KEY(), ""),
	})
	if err != nil {
		return err
	}

	s3Client := s3.New(sess)

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
	return nil
}

func S3Get(counter int) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("sa-east-1"),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(config.AWS_KEY(), config.AWS_SECRET_KEY(), ""),
	})
	if err != nil {
		return "", err
	}

	s3Client := s3.New(sess)

	params := &s3.GetObjectInput{
		Bucket: aws.String(config.AWS_BUCKET()),
		Key:    aws.String(strconv.Itoa(counter) + "-P24_HGUV6.5.3B01_UPGRADE.bin"),
	}

	req, _ := s3Client.GetObjectRequest(params)
	objectURL, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Fatal("[AWS GET LINK]:", params, err)
	}
	return objectURL, nil
}