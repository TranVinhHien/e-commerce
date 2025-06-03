package db

import (
	"context"
	"strings"
)

func buildUpdateProductStockSKU(data []CreateOrderDetailParams, isPlus bool) (string, []interface{}) {
	baseQuery := "UPDATE product_skus SET sku_stock = CASE product_sku_id  "
	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*5)

	for _, vl := range data {
		query := "WHEN ? THEN sku_stock - ?"
		if isPlus {
			query = "WHEN ? THEN sku_stock + ?"
		}
		valueStrings = append(valueStrings, query)
		valueArgs = append(valueArgs, vl.ProductSkuID, vl.Quantity)
	}
	valueStrings_2 := make([]string, 0, len(data))
	for _, vl := range data {
		valueStrings_2 = append(valueStrings_2, "?")
		valueArgs = append(valueArgs, vl.ProductSkuID)
	}
	// query := baseQuery + strings.Join(valueStrings, ",") + " END,update_date = NOW() WHERE product_sku_id IN (" + strings.Join(valueStrings_2, ",") + ");"
	query := baseQuery + strings.Join(valueStrings, " ") + " END, update_date = NOW() WHERE product_sku_id IN (" + strings.Join(valueStrings_2, ",") + ");"
	return query, valueArgs
}
func (s *Queries) UpdateProductStockSKU(ctx context.Context, orderDT []CreateOrderDetailParams, isPlus bool) error {
	query, data := buildUpdateProductStockSKU(orderDT, isPlus)
	_, err := s.db.ExecContext(ctx, query, data...)
	if err != nil {
		return err
	}
	return nil
}

func buildInsertDetailProduct(data []CreateOrderDetailParams) (string, []interface{}) {
	baseQuery := "INSERT INTO order_detail (order_detail_id, quantity, unit_price, product_sku_id, order_id) VALUES "
	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*5)

	for _, vl := range data {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, vl.OrderDetailID, vl.Quantity, vl.UnitPrice, vl.ProductSkuID, vl.OrderID)
	}

	query := baseQuery + strings.Join(valueStrings, ",")
	return query, valueArgs
}
func (s *Queries) InsertDetailProduct(ctx context.Context, orderDT []CreateOrderDetailParams) error {
	query, data := buildInsertDetailProduct(orderDT)
	_, err := s.db.ExecContext(ctx, query, data...)
	if err != nil {
		return err
	}
	return nil
}
