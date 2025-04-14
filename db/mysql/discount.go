package db

import (
	"context"
	"fmt"
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
	rows, err := listData(ctx, s.connPool, table_text, query)
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
	count_sql := fmt.Sprintf("SELECT COUNT(*) as totalElements FROM %s", table_text)

	row := s.connPool.QueryRowContext(ctx, count_sql)
	var sc int64
	err = row.Scan(&sc)
	totalElements = int(sc)
	if err != nil {
		return nil, -1, -1, err
	}
	pages := (float64(totalElements) - 1) / float64(query.PageSize)
	totalPages = int(math.Ceil(pages))

	return
}
