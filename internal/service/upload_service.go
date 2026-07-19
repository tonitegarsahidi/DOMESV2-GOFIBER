package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FileUploadService interface {
	UploadFile(file *multipart.FileHeader, fileType string) (string, error)
}

type fileUploadService struct {
	uploadDir string
}

func NewFileUploadService() FileUploadService {
	dir := "./public/upload"
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Warning: Failed to create upload directory %s: %v\n", dir, err)
	}
	return &fileUploadService{
		uploadDir: dir,
	}
}

func (s *fileUploadService) UploadFile(file *multipart.FileHeader, fileType string) (string, error) {
	if file == nil {
		return "", errors.New("file is nil")
	}

	// Normalize fileType
	fileType = strings.ToLower(strings.TrimSpace(fileType))
	if fileType == "" {
		fileType = "additional" // Default type
	}

	// 1. Validate file extension (Images and Documents only)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isAllowedExtension(ext, fileType) {
		return "", fmt.Errorf("file extension %s is not allowed for type %s", ext, fileType)
	}

	// 2. Validate maximum file size
	maxSize := getMaxSizeForType(fileType)
	if file.Size > maxSize {
		return "", fmt.Errorf("file size exceeds maximum limit of %s", formatBytes(maxSize))
	}

	// 3. Format filename: TYPE-uuid.extension
	// Ensure the prefix matches expected types: cover, main, additional
	prefix := "additional"
	if fileType == "main" || fileType == "primary" {
		prefix = "main"
	} else if fileType == "cover" {
		prefix = "cover"
	}

	newFilename := fmt.Sprintf("%s-%s%s", prefix, uuid.New().String(), ext)
	targetPath := filepath.Join(s.uploadDir, newFilename)

	// 4. Save file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(targetPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// Return public URL path
	return "/public/upload/" + newFilename, nil
}

func isAllowedExtension(ext string, fileType string) bool {
	// Remove leading dot for checking
	e := strings.TrimPrefix(ext, ".")
	
	switch fileType {
	case "cover":
		// Only image files for cover
		allowed := []string{"jpg", "jpeg", "png", "webp"}
		for _, a := range allowed {
			if e == a {
				return true
			}
		}
		return false
	case "main", "primary":
		// Main documents
		allowed := []string{"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx"}
		for _, a := range allowed {
			if e == a {
				return true
			}
		}
		return false
	default:
		// Supporting/additional files (both docs and images allowed)
		allowed := []string{"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "csv", "jpg", "jpeg", "png", "webp", "zip", "rar"}
		for _, a := range allowed {
			if e == a {
				return true
			}
		}
		return false
	}
}

func getMaxSizeForType(fileType string) int64 {
	switch fileType {
	case "cover":
		return 10 * 1024 * 1024 // 10MB
	case "main", "primary":
		return 50 * 1024 * 1024 // 50MB
	default:
		return 50 * 1024 * 1024 // 50MB
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
