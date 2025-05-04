package db

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	db "new-project/db/sqlc"
	services "new-project/services/entity"
	"sync"
)

const (
	payment_online  = "550e8400-e29b-41d4-a716-446655440050"
	payment_offline = "550e8400-e29b-41d4-a716-446655440051"
)

func (s *SQLStore) TXCreateOrdder(ctx context.Context, order *services.Orders, orderDetail []services.OrderDetail) (err error) {
	return s.execTS(ctx, func(tx *db.Queries) error {
		orderDT := make([]db.CreateOrderDetailParams, 0, len(orderDetail))

		orderValue := db.CreateOrderParams{
			OrderID:           order.OrderID,
			TotalAmount:       order.TotalAmount,
			CustomerAddressID: order.CustomerAddressID,
			DiscountID:        sql.NullString{String: order.DiscountID.Data, Valid: order.DiscountID.Valid},
			PaymentMethodID:   order.PaymentMethodID,
			CustomerID:        order.CustomerID,
			OrderStatus:       db.NullOrdersOrderStatus{OrdersOrderStatus: db.OrdersOrderStatus(order.OrderStatus), Valid: order.OrderStatus != ""},
			PaymentStatus:     db.NullOrdersPaymentStatus{OrdersPaymentStatus: db.OrdersPaymentStatus(order.PaymentStatus), Valid: order.PaymentStatus != ""},
		}
		for _, od := range orderDetail {
			orderDT = append(orderDT, db.CreateOrderDetailParams{
				OrderDetailID: od.OrderDetailID,
				Quantity:      od.Quantity,
				UnitPrice:     od.UnitPrice,
				ProductSkuID:  od.ProductSkuID,
				OrderID:       od.OrderID,
			})
		}
		// tru so luong sku
		err := tx.UpdateProductStockSKU(ctx, orderDT)
		if err != nil {
			return err
		}
		// check neu co discount thi tru so luong discount ra
		if orderValue.DiscountID.Valid {
			err = tx.UpdateDiscountAmount(ctx, orderValue.DiscountID.String)
			if err != nil {
				return err
			}
		}
		// tao create order
		err = tx.CreateOrder(ctx, orderValue)
		if err != nil {
			return err
		}
		// tao create orderdetail
		err = tx.InsertDetailProduct(ctx, orderDT)
		if err != nil {
			return err
		}

		return nil
	},
	)
}
func (s *SQLStore) UpdateOrder(ctx context.Context, order services.Orders) (err error) {

	return s.Queries.UpdateOrder(ctx, db.UpdateOrderParams{
		TotalAmount:       sql.NullFloat64{Float64: order.TotalAmount, Valid: order.TotalAmount != 0},
		CustomerAddressID: sql.NullString{String: order.CustomerAddressID, Valid: order.CustomerAddressID != ""},
		DiscountID:        sql.NullString{String: order.DiscountID.Data, Valid: order.DiscountID.Valid},
		PaymentMethodID:   sql.NullString{String: order.PaymentMethodID, Valid: order.PaymentMethodID != ""},
		PaymentStatus:     db.NullOrdersPaymentStatus{OrdersPaymentStatus: db.OrdersPaymentStatus(order.PaymentStatus), Valid: order.PaymentStatus != ""},
		OrderStatus:       db.NullOrdersOrderStatus{OrdersOrderStatus: db.OrdersOrderStatus(order.OrderStatus), Valid: order.OrderStatus != ""},
		OrderID:           order.OrderID,
	})
}
func (s *SQLStore) GetOrderByID(ctx context.Context, orderID string) (i services.Orders, err error) {

	item, err := s.Queries.GetOrder(ctx, orderID)
	if err != nil {
		return i, err
	}
	return item.Convert(), nil
}
func (s *SQLStore) GetOrdersByUserID(ctx context.Context, userID string, query services.QueryFilter) (items []services.Orders, totalPages, totalElements int, err error) {
	table_text := "orders"
	query.Conditions = append(query.Conditions, services.Condition{
		Field:    "customer_id",
		Operator: "=",
		Value:    userID,
	})
	rows, totalElements, err := listData(ctx, s.connPool, table_text, query)
	if err != nil {
		return nil, -1, -1, fmt.Errorf("error listData: %s", err.Error())
	}
	var is []db.Orders
	for rows.Next() {
		var i db.Orders
		if err := rows.Scan(
			&i.OrderID,
			&i.TotalAmount,
			&i.CustomerAddressID,
			&i.DiscountID,
			&i.PaymentMethodID,
			&i.PaymentStatus,
			&i.OrderStatus,
			&i.CreateDate,
			&i.UpdateDate,
			&i.CustomerID,
		); err != nil {
			return nil, -1, -1, fmt.Errorf("error rows.Scan: %s", err.Error())
		}
		is = append(is, i)
	}
	defer rows.Close()
	items = make([]services.Orders, len(is))

	pages := math.Max((float64(totalElements)-1)/float64(query.PageSize), 1)
	totalPages = int(math.Ceil(pages))

	// if len(is)==0 return
	if len(is) == 0 {
		return items, 0, 0, nil
	}
	// get detail and product:
	orderID := make([]string, 0, len(is))
	AddressID := make([]string, 0, len(is))
	addressSet := make(map[string]struct{})
	for _, o := range is {
		orderID = append(orderID, o.OrderID)
		// Kiểm tra nếu o.CustomerAddressID chưa có trong addressSet
		if _, exists := addressSet[o.CustomerAddressID]; !exists {
			AddressID = append(AddressID, o.CustomerAddressID)
			addressSet[o.CustomerAddressID] = struct{}{}
		}
	}
	addressUser, err := s.GetCustomerAddresss(ctx, AddressID)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error GetOrderDetailByorderIDs:%s", err.Error())
	}

	addressMap := make(map[string]services.CustomerAddress, len(addressUser))
	for _, addr := range addressUser {
		addressMap[addr.IDAddress] = addr
	}

	orderDetails, err := s.GetOrderDetailByOrderIDs(ctx, orderID)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error GetOrderDetailByorderIDs:%s", err.Error())
	}
	productSKUIds := make([]string, 0, len(orderDetails))
	for _, o := range orderDetails {
		productSKUIds = append(productSKUIds, o.ProductSkuID)
	}
	producutSKUs, err := s.GetProductsBySKUs(ctx, productSKUIds)
	if err != nil {
		fmt.Println("Error GetOrderDetailByorderIDs:", err)
		return nil, 0, 0, fmt.Errorf("error GetOrderDetailByorderIDs:%s", err.Error())
	}
	if err != nil {
		fmt.Println("Error GetOrderDetailByorderIDs:", err)
		return nil, 0, 0, fmt.Errorf("error GetOrderDetailByorderIDs:%s", err.Error())
	}
	// tạo map truy vấn cho nhanh
	productSkuMap := make(map[string]services.ProductSkusDetail, len(producutSKUs))
	for _, sku := range producutSKUs {
		productSkuMap[sku.ProductSkuID] = sku
	}

	// go
	var wg sync.WaitGroup
	for i, order := range is {
		wg.Add(1)
		go func(idx int, ord db.Orders) {
			// func(idx int, ord db.Orders) {
			defer wg.Done()
			items[idx] = processOrder(ord, orderDetails, productSkuMap, addressMap)
		}(i, order)
	}

	wg.Wait()

	return

}
func processOrder(order db.Orders, orderDetails []services.OrderDetail, productSkuMap map[string]services.ProductSkusDetail, addressMAP map[string]services.CustomerAddress) (OrderService services.Orders) {

	filteredOrderDetails := make([]services.OrderDetail, 0)
	for _, detail := range orderDetails {
		if detail.OrderID == order.OrderID {
			if sku, ok := productSkuMap[detail.ProductSkuID]; ok {
				detail.ProductSKU = services.Narg[services.ProductSkusDetail]{
					Valid: ok,
					Data:  sku,
				}
			}
			filteredOrderDetails = append(filteredOrderDetails, detail)
		}
	}
	OrderService = order.Convert()
	if addr, ok := addressMAP[OrderService.CustomerAddressID]; ok {
		OrderService.Address = services.Narg[services.CustomerAddress]{
			Valid: ok,
			Data:  addr,
		}
	}
	OrderService.OrderDetail = filteredOrderDetails
	paymentName := ""
	if order.PaymentMethodID == payment_online {
		paymentName = "Thanh Toán Onlne"
	}
	if order.PaymentMethodID == payment_offline {
		paymentName = "Thanh Toán Onlne"
	}
	// tạm thời cái payment thì gán cứng cho nhanh
	OrderService.PaymentMethod = services.Narg[services.PaymentMethods]{
		Valid: true,
		Data: services.PaymentMethods{
			PaymentMethodID: order.PaymentMethodID,
			Name:            paymentName,
		},
	}
	return
}

func (s *SQLStore) CheckUserOrder(ctx context.Context, userID, products_spu_id string) (check int64, err error) {
	count, err := s.Queries.CheckUserOrder(ctx, db.CheckUserOrderParams{
		CustomerID:    userID,
		ProductsSpuID: products_spu_id,
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}
