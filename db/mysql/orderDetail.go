package db

import (
	"context"
	"fmt"
	"log"
	db "new-project/db/sqlc"
	services "new-project/services/entity"
	"strings"
)

func (s *SQLStore) buildGetOrderDetailByOrderID(ctx context.Context, orderIDs []string) ([]db.OrderDetail, error) {

	const getProductsBySKU = `-- name: GetProductsBySKU :many
	SELECT order_detail_id, quantity, unit_price, product_sku_id, order_id FROM order_detail
	WHERE order_id IN (%s)
	`

	placeholders := make([]string, len(orderIDs))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf(getProductsBySKU, strings.Join(placeholders, ","))

	// Chuyển đổi danh sách thành các tham số cho câu lệnh SQL
	args := make([]interface{}, len(orderIDs))
	for i, id := range orderIDs {
		args[i] = id
	}

	rows, err := s.connPool.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []db.OrderDetail
	for rows.Next() {
		var i db.OrderDetail
		if err := rows.Scan(
			&i.OrderDetailID,
			&i.Quantity,
			&i.UnitPrice,
			&i.ProductSkuID,
			&i.OrderID,
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

func (s *SQLStore) GetOrderDetailByOrderIDs(ctx context.Context, orderIDs []string) (is []services.OrderDetail, err error) {
	items, err := s.buildGetOrderDetailByOrderID(ctx, orderIDs)
	if err != nil {
		log.Fatal("error when get GetOrderDetailByOrderIDs: ", err)
		return nil, err
	}
	is = make([]services.OrderDetail, len(items))

	// Duyệt và chuyển đổi
	for i, item := range items {
		is[i] = services.OrderDetail{
			OrderDetailID: item.OrderDetailID,
			Quantity:      item.Quantity,
			UnitPrice:     item.UnitPrice,
			ProductSkuID:  item.ProductSkuID,
			OrderID:       item.OrderID,
		}
	}
	return
}
