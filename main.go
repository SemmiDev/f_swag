package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	_ "github.com/semmidev/s_swag/docs/f_swag"
)

// @title Upload File API
// @description Upload File API with Go and Fiber
// @version 1

// @host localhost:3000
// @BasePath /api
// @schemes http
func main() {
	app := NewHTTPServer()
	log.Fatal(app.Listen(":3000"))
}

func NewHTTPServer() *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit: 25 * 1024 * 1024,
	})

	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/swagger/*", swagger.HandlerDefault) // default
	app.Post("/api/upload", uploadFile)
	app.Post("/api/upload/multiple", uploadMultipleFiles)

	return app
}

func uploadMultipleFiles(c *fiber.Ctx) error {
	// Menerima semua file dari request
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	// Mendapatkan semua file yang diunggah
	files := form.File["files"]

	// Looping untuk setiap file
	for _, file := range files {
		// Mendapatkan ekstensi file
		ext := filepath.Ext(file.Filename)

		// Generate UUID sebagai nama file baru
		filename := uuid.New().String() + ext

		// Simpan file ke dalam folder
		if err := saveFile(file, filename); err != nil {
			return err
		}
	}

	// Tampilkan response
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": fmt.Sprintf("%d file berhasil diupload.", len(files)),
	})
}

// uploadFile godoc
// @Summary Upload file
// @Description Upload file to server
// @Tags upload
// @Accept mpfd
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]interface{}
// @Router /upload [post]
func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	ext := filepath.Ext(file.Filename)

	fmt.Println("---- File Info ----")
	fmt.Println("File Name: ", file.Filename)
	fmt.Println("File Ext: ", ext)
	fmt.Println("File Size: ", file.Size/1024/1024, "MB")
	fmt.Println("------------------")

	// Generate UUID sebagai nama file baru
	filename := uuid.New().String() + ext

	// Simpan file ke dalam folder
	if err := saveFile(file, filename); err != nil {
		return err
	}

	// Tampilkan response
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": fmt.Sprintf("File %s berhasil diupload.", file.Filename),
	})
}

func saveFile(file *multipart.FileHeader, filename string) error {
	// Membuka file yang diterima
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Membaca isi file
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	// Membuat folder jika belum ada
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	// Menyimpan file ke dalam folder uploads dengan nama baru
	dst, err := os.Create(filepath.Join("uploads", filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.Write(data)
	if err != nil {
		return err
	}

	return nil
}
