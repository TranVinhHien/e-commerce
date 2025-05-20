package db

import (
	"context"
	"fmt"
	"log"
	"math"
	db "new-project/db/sqlc"
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

func (s *SQLStore) buildGetPrice(ctx context.Context, productSpuIDs []string) ([]ProductSkusDetailss, error) {

	const getProductsBySKU = `-- name: GetProductsBySKU :many
	WITH RankedSkus AS (
    SELECT 
        sku.price,
        sku.products_spu_id,
        ROW_NUMBER() OVER (PARTITION BY sku.products_spu_id ORDER BY sku.create_date ASC) AS rn
    FROM product_skus AS sku
    WHERE sku.products_spu_id IN (%s)
	)
	SELECT products_spu_id,price
	FROM RankedSkus
	WHERE rn = 1;
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

	var items []ProductSkusDetailss
	for rows.Next() {
		var i ProductSkusDetailss
		if err := rows.Scan(
			&i.ProductsSpuID,
			&i.Price,
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

func (s *SQLStore) GetAllProductSimple(ctx context.Context, query services.QueryFilter) (items []services.ProductSimple, totalPages, totalElements int, err error) {
	table_text := "products_spu"
	rows, totalElements, err := listData(ctx, s.connPool, table_text, query)
	if err != nil {
		return nil, -1, -1, err
	}
	var is []db.ProductsSpu
	spuIDs := make([]string, 0, len(is))

	for rows.Next() {
		var i db.ProductsSpu
		if err := rows.Scan(
			&i.ProductsSpuID,
			&i.Name,
			&i.BrandID,
			&i.Description,
			&i.ShortDescription,

			&i.StockStatus,
			&i.DeleteStatus,
			&i.Sort,
			&i.CreateDate,
			&i.UpdateDate,
			&i.Image,

			&i.Media,
			&i.Key,
			&i.CategoryID,
		); err != nil {
			return nil, -1, -1, err
		}
		spuIDs = append(spuIDs, i.ProductsSpuID)
		is = append(is, i)
	}
	defer rows.Close()
	spuPriceMap := make(map[string]float64)
	if len(spuIDs) == 0 {
		return nil, -1, -1, fmt.Errorf("not found ")
	}
	// get price
	priceSPUs, err := s.buildGetPrice(ctx, spuIDs)
	for _, o := range priceSPUs {
		spuPriceMap[o.ProductsSpuID] = o.Price
	}
	spuRating := make(map[string]RatingAvgAndCountSPU)

	// get price
	ratings, err := s.buildGetRatingsAVGAndCOUNT(ctx, spuIDs)
	for _, o := range ratings {
		spuRating[o.ProductsSpuID] = o
	}

	items = make([]services.ProductSimple, len(is))

	// Duyệt và chuyển đổi
	for i, item := range is {
		spu := item.Convert()
		item := services.ProductSimple{
			ProductsSpuID:    spu.ProductsSpuID,
			Name:             spu.Name,
			ShortDescription: spu.ShortDescription,
			StockStatus:      spu.StockStatus,
			DeleteStatus:     spu.DeleteStatus,
			Image:            spu.Image,
			// Media:            spu.Media,
			Avg_star:    spuRating[spu.ProductsSpuID].Avg_star,
			TotalRating: spuRating[spu.ProductsSpuID].TotalRating,
			Key:         spu.Key,
			CategoryID:  spu.CategoryID,
			Price:       spuPriceMap[spu.ProductsSpuID],
		}
		items[i] = item
	}

	pages := (float64(totalElements) - 1) / float64(query.PageSize)
	totalPages = int(math.Ceil(pages))

	return
}
func (s *SQLStore) GetProductDetail(ctx context.Context, productSpuID string) (product_detail services.ProductDetail, err error) {

	spu, err := s.Queries.GetProductSPU(ctx, productSpuID)
	if err != nil {
		return
	}
	skus, err := s.Queries.ListProductSKUs(ctx, productSpuID)
	if err != nil {
		return
	}
	des_attr, err := s.Queries.ListDescriptionAttrs(ctx, productSpuID)
	if err != nil {
		return
	}
	sku_attr, err := s.Queries.ListProductSKUAttrs(ctx, productSpuID)
	if err != nil {
		return
	}
	sku := make([]services.ProductSkus, 0, len(skus))
	desAttr := make([]services.DescriptionAttr, 0, len(des_attr))
	skuAttr := make([]services.ProductSkuAttrs, 0, len(sku_attr))
	for _, o := range skus {
		sku = append(sku, o.Convert())
	}
	for _, o := range des_attr {
		desAttr = append(desAttr, o.Convert())
	}
	for _, o := range sku_attr {
		skuAttr = append(skuAttr, o.Convert())
	}
	product_detail.Spu = spu.Convert()
	product_detail.Sku = sku
	product_detail.DesAttr = desAttr
	product_detail.SkuAttr = skuAttr
	return
}
