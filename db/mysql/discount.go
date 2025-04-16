package db

import (
	"context"
	"math"
	db "new-project/db/sqlc"
	services "new-project/services/entity"
)

func (s *SQLStore) Discount(ctx context.Context, discount string) (i services.Discounts, err error) {
	row, err := s.Queries.GetDiscount(ctx, discount)
	if err != nil {
		return services.Discounts{}, err
	}
	return row.Convert(), nil
}

func (s *SQLStore) ListDiscount(ctx context.Context, query services.QueryFilter) (items []services.Discounts, totalPages, totalElements int, err error) {
	table_text := "discounts"
	rows, totalElements, err := listData(ctx, s.connPool, table_text, query)
	if err != nil {
		return nil, -1, -1, err
	}
	var is []db.Discounts
	for rows.Next() {
		var i db.Discounts
		if err := rows.Scan(
			&i.DiscountID,
			&i.DiscountCode,
			&i.DiscountValue,
			&i.StartDate,
			&i.EndDate,
			&i.MinOrderValue,
			&i.Amount,
			&i.CreateDate,
			&i.UpdateDate,
		); err != nil {
			return nil, -1, -1, err
		}
		is = append(is, i)
	}
	defer rows.Close()
	items = make([]services.Discounts, len(is))

	// Duyệt và chuyển đổi
	for i, item := range is {
		items[i] = item.Convert()
	}

	pages := (float64(totalElements) - 1) / float64(query.PageSize)
	totalPages = int(math.Ceil(pages))

	return
}
