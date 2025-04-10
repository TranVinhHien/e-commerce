package db

import (
	"context"
	"log"
	services "new-project/services/entity"
)

func (s *SQLStore) ListCategories(ctx context.Context) ([]services.Categorys, error) {
	is, err := s.Queries.ListCategories(ctx)
	if err != nil {
		log.Fatal("error when get ListCategories: ", err)
		return nil, err
	}
	items := make([]services.Categorys, len(is))

	// Duyệt và chuyển đổi
	for i, item := range is {
		items[i] = item.Convert()
	}
	return items, nil
}
func (s *SQLStore) ListCategoriesByID(ctx context.Context, cate_id string) ([]services.Categorys, error) {
	is, err := s.Queries.ListCategoriesByID(ctx, cate_id)
	if err != nil {
		log.Fatal("error when get ListCategoriesByID: ", err)
		return nil, err
	}
	items := make([]services.Categorys, len(is))

	// Duyệt và chuyển đổi
	for i, item := range is {
		items[i] = item.Convert()
	}
	return items, nil
}
