package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func handleGetFileName(c *fiber.Ctx) error {
	fileNameParam := c.Params("filename")
	filePath := path.Join("storage", "product", fileNameParam)

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return c.Status(http.StatusNotFound).SendString("File not found")
		}

		log.Println(err)
		return c.Status(http.StatusInternalServerError).SendString("Terjadi kesalahan saat mengakses file")
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).SendString("Terjadi kesalahan saat membuka file")
	}

	ext := path.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)

	c.Set("Content-Type", mimeType)
	return c.SendStream(file)
}

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Get("/storage/product/:filename", handleGetFileName)
	app.Post("/product/upload", handler.UploadProductImageHandler)

	app.Listen(":3000")
}
