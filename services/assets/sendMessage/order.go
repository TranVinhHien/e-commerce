package services_assets_sendMessage

import (
	"firebase.google.com/go/messaging"
)

func ThanhToanThanhCong() *messaging.Notification {
	return &messaging.Notification{
		Title:    "Thanh toán thành công",
		Body:     "Cảm ơn bạn đã thanh toán đơn hàng.",
		ImageURL: "https://f6e9-118-68-56-216.ngrok-free.app/v1/media/products?id=images/phu-kien-thoi-trang___phu-kien-nu___phu-kien-nu-khac/100477499_1.jpg",
	}
}
