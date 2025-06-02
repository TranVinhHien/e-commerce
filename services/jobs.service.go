package services

import (
	"context"
)

func (s *service) NotiNewDiscount(ctx context.Context) error {
	// query.Conditions = []services.Condition{
	// 	{Field: "amount", Operator: ">=", Value: 1},
	// 	{Field: "end_date", Operator: ">", Value: time.Now()},
	// }
	// discounts, totalPages, totalElements, err := s.repository.ListDiscount(ctx, query)

	// if err != nil {
	// 	fmt.Println("Error ListDiscount:", err)
	// 	return nil, assets_services.NewError(400, err)
	// }

	// result, err := assets_services.HideFields(discounts, "discounts")
	// if err != nil {
	// 	fmt.Println("Error HideFields:", err)
	// 	return nil, assets_services.NewError(400, err)
	// }
	// result["currentPage"] = query.Page
	// result["totalPages"] = totalPages
	// result["totalElements"] = totalElements
	// result["limit"] = query.PageSize
	// return result, nil
	return nil
}
