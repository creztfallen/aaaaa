package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// HttpUploadFileHandler is a util function to grab files from post request forms using net/http
func HttpUploadFileHandler(w http.ResponseWriter, r *http.Request) {

	// check if the used method is POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// grab the file from the form
	file, header, err := r.FormFile("filetoupload")
	if err != nil {
		log.Println("Couldn't get the file:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = file.Close()
	if err != nil {
		log.Println("Couldn't close the file:", err)
		return
	}

	// in case the directory to put the file in doesn't exist, create one
	err = os.MkdirAll("config-file-downloads", os.ModePerm)
	if err != nil {
		log.Println("Couldn't create the directory 'config-file-downloads':", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create the whole path to save the file
	filename := header.Filename
	uploadFilepath := filepath.Join("config-file-downloads", filename)

	// create the file
	outFile, err := os.Create(uploadFilepath)
	if err != nil {
		log.Println("Couldn't create the file:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = outFile.Close()
	if err != nil {
		log.Println("Couldn't close the file:", err)
		return
	}

	// read the file content
	_, err = io.Copy(outFile, file)
	if err != nil {
		log.Println("Couldn't save the file:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("File received and successfully saved"))
	if err != nil {
		log.Println("Couldn't write data to connection:", err)
	}
}

// UploadFileHandler is a util function to grab files from post request forms using fiber
func UploadFileHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("filetoupload")
	if err != nil {
		log.Println("Couldn't get the file:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Error trying to get the file")
	}

	err = c.SaveFile(file, fmt.Sprintf("./cmd/api/compile/config-file-downloads/%s", file.Filename))
	if err != nil {
		log.Println("Couldn't save the file:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error trying to save the file")
	}

	return c.SendString("File received and saved successfully")
}
