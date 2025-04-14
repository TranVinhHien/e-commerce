package db

import (
	"context"
	"log"
	services "new-project/services/entity"
)

func (s *SQLStore) ListPayment(ctx context.Context) (items []services.PaymentMethods, err error) {
	is, err := s.Queries.ListPaymentMethods(ctx)
	if err != nil {
		log.Fatal("error when get ListCategories: ", err)
		return nil, err
	}
	items = make([]services.PaymentMethods, len(is))

	// Duyệt và chuyển đổi
	for i, item := range is {
		items[i] = item.Convert()
	}
	return items, nil
}
