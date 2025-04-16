package db

import (
	"context"
	"fmt"
	"log"
	services "new-project/services/entity"
	"strings"
	"time"
)

func (s *SQLStore) buildGetSKU(ctx context.Context, productSkuIDs []string) ([]ProductSkusDetailss, error) {

	const getProductsBySKU = `-- name: GetProductsBySKU :many
	SELECT sku.product_sku_id, sku.value, sku.sku_stock, sku.price, sku.sort, sku.create_date, sku.update_date, sku.products_spu_id,
	spu.name,spu.short_description,spu.image,
	 (
      SELECT GROUP_CONCAT(CONCAT(attr.name, ':', attr.value) SEPARATOR ', ')
      FROM product_sku_attrs AS attr
      WHERE FIND_IN_SET(attr.product_sku_attr_id, REPLACE(sku.value, '/', ',')) > 0
    ) AS info_sku_attr
	FROM product_skus as sku join products_spu as spu on sku.products_spu_id = spu.products_spu_id
	WHERE sku.product_sku_id IN (%s)
	`

	placeholders := make([]string, len(productSkuIDs))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf(getProductsBySKU, strings.Join(placeholders, ","))

	// Chuyển đổi danh sách thành các tham số cho câu lệnh SQL
	args := make([]interface{}, len(productSkuIDs))
	for i, id := range productSkuIDs {
		args[i] = id
	}

	rows, err := s.connPool.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ProductSkusDetailss
	for rows.Next() {
		var i ProductSkusDetailss
		if err := rows.Scan(
			&i.ProductSkuID,
			&i.Value,
			&i.SkuStock,
			&i.Price,
			&i.Sort,
			&i.CreateDate,
			&i.UpdateDate,
			&i.ProductsSpuID,
			&i.Name,
			&i.ShortDescription,
			&i.Image,
			&i.InfoProduct,
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

func (s *SQLStore) GetProductsBySKUs(ctx context.Context, product_sku_ids []string) (is []services.ProductSkusDetail, err error) {
	items, err := s.buildGetSKU(ctx, product_sku_ids)
	if err != nil {
		log.Fatal("error when get GetProductsBySKUs: ", err)
		return nil, err
	}
	is = make([]services.ProductSkusDetail, len(items))

	// Duyệt và chuyển đổi
	for i, item := range items {
		is[i] = services.ProductSkusDetail{
			ProductSkuID:     item.ProductSkuID,
			Value:            item.Value,
			SkuStock:         item.SkuStock.Int32,
			Price:            item.Price,
			Sort:             item.Sort.Int32,
			CreateDate:       item.CreateDate.Time,
			UpdateDate:       services.Narg[time.Time]{Data: item.UpdateDate.Time, Valid: item.UpdateDate.Valid},
			ProductsSpuID:    item.ProductsSpuID,
			Name:             item.Name,
			ShortDescription: item.ShortDescription,
			Image:            item.Image,
			InfoProduct:      item.InfoProduct,
		}
	}
	return
}
