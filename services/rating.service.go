package services

import (
	"context"
	"fmt"
	assets_services "new-project/services/assets"
	services "new-project/services/entity"

	"github.com/google/uuid"
)

func (s *service) ListRating(ctx context.Context, products_spu_id string, query services.QueryFilter) (map[string]interface{}, *assets_services.ServiceError) {
	query.Conditions = []services.Condition{
		{Field: "products_spu_id", Operator: "=", Value: products_spu_id},
		{Field: "comment", Operator: "!=", Value: ""},
	}
	discounts, totalPages, totalElements, err := s.repository.GetRatings(ctx, query)

	if err != nil {
		fmt.Println("Error ListDiscount:", err)
		return nil, assets_services.NewError(400, err)
	}

	result, err := assets_services.HideFields(discounts, "discounts", "account_id", "customer_id")
	if err != nil {
		fmt.Println("Error HideFields:", err)
		return nil, assets_services.NewError(400, err)
	}
	result["currentPage"] = query.Page
	result["totalPages"] = totalPages
	result["totalElements"] = totalElements
	result["limit"] = query.PageSize
	return result, nil
}
func (s *service) CreateRating(ctx context.Context, userID string, rating services.Ratings) (err *assets_services.ServiceError) {

	// check nguoi dung phai mua sanpham moi  duoc binh  luan
	count, errS := s.repository.CheckUserOrder(ctx, userID, rating.ProductsSpuID)
	if errS != nil {
		return assets_services.NewError(400, errS)
	}
	if count == 0 {
		return assets_services.NewError(400, fmt.Errorf("tài khoảng chưa mua sản phẩm không thể dánh giá"))
	}
	errS = s.repository.CreateRating(ctx, services.Ratings{
		RatingID:      uuid.New().String(),
		Comment:       rating.Comment,
		Star:          rating.Star,
		ProductsSpuID: rating.ProductsSpuID,
		CustomerID:    userID,
	})
	if errS != nil {
		return assets_services.NewError(400, errS)
	}
	return nil
}
