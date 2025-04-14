package services

import (
	"context"
	"fmt"
	assets_services "new-project/services/assets"
)

func (s *service) ListPayment(ctx context.Context) (map[string]interface{}, *assets_services.ServiceError) {

	payments, err := s.repository.ListPayment(ctx)
	if err != nil {
		fmt.Println("Error ListPayment:", err)
		return nil, assets_services.NewError(400, err)
	}
	result, err := assets_services.HideFields(payments, "payments")
	if err != nil {
		fmt.Println("Error HideFields:", err)
		return nil, assets_services.NewError(400, err)
	}
	return result, nil
}
