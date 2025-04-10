package assets_services

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"reflect"
)

func HideFields(obj interface{}, key string, fieldsToHide ...string) (map[string]interface{}, error) {
	// Convert object to JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	// Xử lý trường hợp đầu vào là slice/array
	if reflect.TypeOf(obj).Kind() == reflect.Slice || reflect.TypeOf(obj).Kind() == reflect.Array {
		var slice []map[string]interface{}
		if err := json.Unmarshal(data, &slice); err != nil {
			return nil, err
		}

		// Ẩn các field chỉ định trong từng phần tử
		if fieldsToHide != nil {
			for i := range slice {
				for _, field := range fieldsToHide {
					delete(slice[i], field)
				}
			}
		}

		// Trả về dạng map với key "data" chứa mảng
		return map[string]interface{}{
			key: slice,
		}, nil
	}

	// Xử lý trường hợp đầu vào là đối tượng đơn
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

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
