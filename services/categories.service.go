package services

import (
	"context"
	"fmt"
	assets_services "new-project/services/assets"
)

func (s *service) GetCategoris(ctx context.Context, userName string) (map[string]interface{}, *assets_services.ServiceError) {

	caterd, err := s.redis.GetCategoryTree(ctx, userName)
	if err != nil {
		fmt.Println("Error GetCategories:", err)
		return nil, assets_services.NewError(400, err)
	}
	if len(caterd) == 0 {
		cates, err := s.repository.ListCategories(ctx)
		if err != nil {
			fmt.Println("Error ListCategories:", err)
			return nil, assets_services.NewError(400, err)
		}
		err = s.redis.AddCategories(ctx, cates)
		if err != nil {
			fmt.Println("Error AddCategories:", err)
			return nil, assets_services.NewError(400, err)
		}
		caterd, err = s.redis.GetCategoryTree(ctx, userName)
		if err != nil {
			fmt.Println("Error GetCategories:", err)
			return nil, assets_services.NewError(400, err)
		}
	}
	// fmt.Println("\nGet all categories:")
	// for _, cat := range caterd {
	// 	fmt.Printf("Category: %s (%s)\n", cat.Name, cat.CategoryID)
	// 	if cat.Childs.Valid {
	// 		for _, child := range cat.Childs.Data {
	// 			fmt.Printf("  - Child: %s (%s)\n", child.Name, child.CategoryID)
	// 			if child.Childs.Valid {
	// 				for _, child2 := range child.Childs.Data {
	// 					fmt.Printf("	  - Child 2: %s (%s)\n", child2.Name, child2.CategoryID)

	// 				}
	// 			}
	// 		}
	// 	}
	// }

	result, err := assets_services.HideFields(caterd, "categories", "parent")
	if err != nil {
		fmt.Println("Error HideFields:", err)
		return nil, assets_services.NewError(400, err)
	}
	return result, nil
}
