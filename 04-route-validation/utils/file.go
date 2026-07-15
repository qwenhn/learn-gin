package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var allowExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

var allowMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

const maxSize = 5 << 20 // (5MB)

func ValidateAndSaveFile(fileHeader *multipart.FileHeader, path string) (string, error) {
	// check extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowExts[ext] {
		return "", errors.New("unsupported file extension")
	}

	if fileHeader.Size > maxSize {
		return "", errors.New("File too large (5MB)")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", errors.New("Unable to open the file")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", errors.New("Unable to read the file")
	}

	mimeType := http.DetectContentType(buffer)
	if !allowMimeTypes[mimeType] {
		return "", fmt.Errorf("invalid MIME type: %s", mimeType)
	}

	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		return "", errors.New("Unable to create the upload folder")
	}

	savePath := filepath.Join(path, filename)
	if err := saveFile(fileHeader, savePath); err != nil {
		return "", err
	}

	return filename, nil
}

func saveFile(fileHeader *multipart.FileHeader, destination string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)

	return err
}
