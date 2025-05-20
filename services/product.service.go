package services

import (
	"context"
	"fmt"
	assets_services "new-project/services/assets"
	services "new-project/services/entity"
)

func (s *service) GetAllProductSimple(ctx context.Context, query services.QueryFilter) (map[string]interface{}, *assets_services.ServiceError) {
	product_spu, totalPages, totalElements, err := s.repository.GetAllProductSimple(ctx, query)
	if err != nil {
		fmt.Println("Error GetAllProductSimple:", err)
		return nil, assets_services.NewError(400, err)
	}

	result, err := assets_services.HideFields(product_spu, "products")
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
func (s *service) GetDetailProduct(ctx context.Context, productSpuID string) (map[string]interface{}, *assets_services.ServiceError) {
	product_spi_detail, err := s.repository.GetProductDetail(ctx, productSpuID)
	if err != nil {
		fmt.Println("Error GetProductDetail:", err)
		return nil, assets_services.NewError(400, err)
	}

	result, err := assets_services.HideFields(product_spi_detail, "product")
	if err != nil {
		fmt.Println("Error HideFields:", err)
		return nil, assets_services.NewError(400, err)
	}

	return result, nil
}
