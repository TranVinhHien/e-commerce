package db

import (
	"context"
	"database/sql"
	"math"
	db "new-project/db/sqlc"
	services "new-project/services/entity"
)

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
