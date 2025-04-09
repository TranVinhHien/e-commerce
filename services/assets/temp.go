package assets_services

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

func HideFields(obj interface{}, fieldsToHide ...string) (map[string]interface{}, error) {
	// Convert object to JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON to map
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	// Remove specified fields
	for _, field := range fieldsToHide {
		delete(result, field)
	}
	return result, nil
}
func SaveUploadedFile(fileHeader *multipart.FileHeader, destination string) error {
	// Mở file gốc từ form
	src, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("mở file lỗi: %v", err)
	}
	defer src.Close()

	// Tạo file đích để lưu
	dst, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("không thể tạo file: %v", err)
	}
	defer dst.Close()

	// Ghi nội dung từ src -> dst
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("ghi file lỗi: %v", err)
	}

	return nil
}
