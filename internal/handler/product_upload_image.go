package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UploadProductImageHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data gambar tidak ditemukan",
		})
	}

	// Validasi ekstensi file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Ekstensi file tidak valid. Hanya diperbolehkan: .jpg, .jpeg, .png, .webp",
		})
	}

	// validasi content type
	contentType := file.Header.Get("Content-Type")
	allowedContentType := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}

	if !allowedContentType[contentType] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Content type tidak valid. Hanya diperbolehkan: image/jpeg, image/png, image/webp",
		})
	}

	timeStamp := time.Now().UnixNano()
	fileName := fmt.Sprintf("product_%d%s", timeStamp, filepath.Ext(file.Filename))
	uploadPath := "./storage/product/" + fileName
	err = c.SaveFile(file, uploadPath)
	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"message":   "Upload gambar produk berhasil",
		"file_name": fileName,
	})
}
