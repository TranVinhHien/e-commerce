package services

import (
	"context"
)

func (s *service) RenderImage(ctx context.Context, filename string) string {
	// kiểm tra thông tin abc xyz'
	filePath := s.env.ImagePath + filename
	return filePath
}
