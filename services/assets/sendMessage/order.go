package services_assets_sendMessage

import (
	"fmt"

	"firebase.google.com/go/messaging"
)

func ThanhToanThanhCong(orderID string, total float64) *messaging.Notification {
	return &messaging.Notification{
		Title: "Thanh toán thành công",
		Body:  "Bạn đã thanh toán thành công đơn hàng." + orderID + "với giá trị" + fmt.Sprint("%f", total) + "VNĐ",
		// ImageURL: "https://f6e9-118-68-56-216.ngrok-free.app/v1/media/products?id=images/phu-kien-thoi-trang___phu-kien-nu___phu-kien-nu-khac/100477499_1.jpg",
	}
}
