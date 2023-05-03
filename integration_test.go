package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUploadMultipleFiles(t *testing.T) {
	// Membuat buffer untuk menampung body request
	body := &bytes.Buffer{}

	// Membuat multipart writer untuk menulis form-data
	writer := multipart.NewWriter(body)

	// Membuka file yang akan diupload
	file1, err := os.Open("uploads/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file1.Close()

	// Mendapatkan informasi file
	stat, err := file1.Stat()
	if err != nil {
		t.Fatal(err)
	}

	// Membuat bagian form-data untuk file pertama
	part1, err := writer.CreateFormFile("files", stat.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Menyalin isi file ke dalam bagian form-data
	if _, err := io.Copy(part1, file1); err != nil {
		t.Fatal(err)
	}

	// Membuka file kedua yang akan diupload
	// file2, err := os.Open("uploads/test.txt")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer file2.Close()

	// // Mendapatkan informasi file
	// stat, err = file2.Stat()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // Membuat bagian form-data untuk file kedua
	// part2, err := writer.CreateFormFile("files", stat.Name())
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // Menyalin isi file ke dalam bagian form-data
	// if _, err := io.Copy(part2, file2); err != nil {
	// 	t.Fatal(err)
	// }

	// Menyelesaikan penulisan form-data
	writer.Close()

	// Membuat HTTP request dengan method POST dan body form-data
	req := httptest.NewRequest(http.MethodPost, "/api/upload/multiple", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Membuat HTTP handler dari Fiber app
	app := NewHTTPServer()

	response, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	// Memastikan status code yang diharapkan sesuai dengan status code yang didapatkan
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Memastikan response body yang diharapkan sesuai dengan response body yang didapatkan
	expected := `{"status":"success"}`
	responseBody, _ := io.ReadAll(response.Body)

	responseBodyMap := map[string]interface{}{}
	if err := json.Unmarshal(responseBody, &responseBodyMap); err != nil {
		t.Fatal(err)
	}

	if responseBodyMap["status"] != "success" {
		t.Fatalf("Expected response body %s, got %s", expected, responseBody)
	}
}
