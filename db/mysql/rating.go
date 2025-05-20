package db

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	db "new-project/db/sqlc"
	services "new-project/services/entity"
	"strings"
)

type RatingAvgAndCountSPU struct {
	ProductsSpuID string  `json:"products_spu_id"`
	Avg_star      float64 `json:"average_star"`
	TotalRating   int32   `json:"total_rating"`
}

func (s *SQLStore) buildGetRatingsAVGAndCOUNT(ctx context.Context, productSpuIDs []string) ([]RatingAvgAndCountSPU, error) {

	const getProductsBySKU = `
	SELECT 
    r.products_spu_id, 
    COUNT(r.rating_id) AS total_rating, 
    ROUND(AVG(r.star), 1)  AS average_star
FROM 
    ratings r
WHERE 
    r.products_spu_id IN (%s)  -- Thay thế với danh sách ID của sản phẩm
GROUP BY 
    r.products_spu_id;
	`

	placeholders := make([]string, len(productSpuIDs))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf(getProductsBySKU, strings.Join(placeholders, ","))

	// Chuyển đổi danh sách thành các tham số cho câu lệnh SQL
	args := make([]interface{}, len(productSpuIDs))
	for i, id := range productSpuIDs {
		args[i] = id
	}

	rows, err := s.connPool.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []RatingAvgAndCountSPU
	for rows.Next() {
		var i RatingAvgAndCountSPU
		if err := rows.Scan(
			&i.ProductsSpuID,
			&i.TotalRating,
			&i.Avg_star,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *SQLStore) GetRatings(ctx context.Context, query services.QueryFilter) (items []services.Ratings, totalPages, totalElements int, err error) {
	table_text := "ratings"

	rows, totalElements, err := listData(ctx, s.connPool, table_text, query)
	if err != nil {
		return nil, -1, -1, err
	}
	var is []db.Ratings
	for rows.Next() {
		var i db.Ratings
		if err := rows.Scan(
			&i.RatingID,
			&i.Comment,
			&i.Star,
			&i.CreateDate,
			&i.UpdateDate,
			&i.CustomerID,
			&i.ProductsSpuID,
		); err != nil {
			return nil, -1, -1, err
		}
		is = append(is, i)
	}
	defer rows.Close()
	items = make([]services.Ratings, len(is))
	// Duyệt và chuyển đổi
	for i, item := range is {
		items[i] = item.Convert()
	}

	pages := (float64(totalElements) - 1) / float64(query.PageSize)
	totalPages = int(math.Ceil(pages))

	return
}
func (s *SQLStore) CreateRating(ctx context.Context, rating services.Ratings) (err error) {
	ratingPR := db.CreateRatingParams{
		RatingID: rating.RatingID,
		Comment: sql.NullString{
			String: rating.Comment.Data,
			Valid:  rating.Comment.Valid,
		},
		Star:          rating.Star,
		CustomerID:    rating.CustomerID,
		ProductsSpuID: rating.ProductsSpuID,
	}
	err = s.Queries.CreateRating(ctx, ratingPR)
	if err != nil {
		return err
	}
	return nil
}
