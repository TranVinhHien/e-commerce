package db

import (
	"context"
	"database/sql"
	db "new-project/db/sqlc"
	services "new-project/services/entity"
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
