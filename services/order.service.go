package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	config_assets "new-project/assets/config"
	assets_services "new-project/services/assets"
	services_assets_sendMessage "new-project/services/assets/sendMessage"
	services "new-project/services/entity"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func createMoMoPayload(env config_assets.ReadENV, payloadd services.CombinedDataPayLoadMoMo) services.Payload_MOMO {

	var partnerCode = "MOMO"
	var extraData = ""
	var partnerName = "Test"
	var storeId = "MoMoTestStore"
	var orderGroupId = ""
	var autoCapture = true
	var lang = "vi"
	var requestType = "payWithMethod"
	var amount = strconv.Itoa(int(payloadd.OrderTX.TotalAmount))
	var orderId = uuid.New().String()
	var orderInfo = fmt.Sprintf("Thanh toán %s VNĐ cho đơn hàng : %s", amount, orderId)
	var accessKey = env.AccessKeyMoMo
	var secretKey = env.SecretKeyMoMo
	var redirectUrl = env.RedirectURL
	var ipnUrl = env.IpnURL
	var requestId = payloadd.OrderTX.OrderID

	var rawSignature bytes.Buffer
	rawSignature.WriteString("accessKey=")
	rawSignature.WriteString(accessKey)
	rawSignature.WriteString("&amount=")
	rawSignature.WriteString(amount)
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString(extraData)
	rawSignature.WriteString("&ipnUrl=")
	rawSignature.WriteString(ipnUrl)
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(orderId)
	rawSignature.WriteString("&orderInfo=")
	rawSignature.WriteString(orderInfo)
	rawSignature.WriteString("&partnerCode=")
	rawSignature.WriteString(partnerCode)
	// chuyển hướng về ứng dụng của mình
	rawSignature.WriteString("&redirectUrl=")
	rawSignature.WriteString(redirectUrl)
	rawSignature.WriteString("&requestId=")
	rawSignature.WriteString(requestId)
	rawSignature.WriteString("&requestType=")
	rawSignature.WriteString(requestType)

	// Create a new HMAC by defining the hash type and the key (as byte array)
	hmac := hmac.New(sha256.New, []byte(secretKey))

	// Write Data to it
	hmac.Write(rawSignature.Bytes())

	// Get result and encode as hexadecimal string
	signature := hex.EncodeToString(hmac.Sum(nil))

	var payload = services.Payload_MOMO{
		PartnerCode:  partnerCode,
		AccessKey:    accessKey,
		RequestID:    requestId,
		Amount:       amount,
		RequestType:  requestType,
		RedirectUrl:  redirectUrl,
		IpnUrl:       ipnUrl,
		OrderID:      orderId,
		StoreId:      storeId,
		PartnerName:  partnerName,
		OrderGroupId: orderGroupId,
		AutoCapture:  autoCapture,
		Lang:         lang,
		OrderInfo:    orderInfo,
		ExtraData:    extraData,
		Signature:    signature,
		Items:        payloadd.Items,
		UserInfo: services.User_MOMO{
			Name:        payloadd.Info.Name,
			PhoneNumber: payloadd.Address.PhoneNumber,
			Email:       payloadd.Address.Address,
		},
	}

	return payload
}

func (s *service) CreateOrder(ctx context.Context, user_id string, order *services.CreateOrderParams) (map[string]interface{}, *assets_services.ServiceError) {
	// check neu co order thi khong cho nguoi dung thanh toan
	if vl, _ := s.redis.GetOrderOnline(ctx, user_id); vl != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("bạn còn order:  %s chưa được thanh toán", vl.OrderTX.OrderID))
	}
	// lay thong tin dia chi nguoi dung va kiem tra hop le
	info, address, err := s.repository.CustomerAddresses(ctx, user_id, order.Address_id)
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("CustomerAddresses %s", err.Error()))
	}
	// lay thong tin tat ca san pham  order.NumOfProducts[0].Product_sku_id
	product_sku_ids := make([]string, 0)
	for _, item := range order.NumOfProducts {
		product_sku_ids = append(product_sku_ids, item.Product_sku_id)
	}
	product_sku, err := s.repository.GetProductsBySKUs(ctx, product_sku_ids)
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("GetProductsBySKUs %s", err.Error()))
	}
	// kiem tra san pham co con hang hay khong
	//	// neu co 1 san pham het hang thì  trả về lỗi cho người dùng

	// thực hiện tính tổng tiền cho từng loại sản phẩm và tổng tiền của order
	// thực hiện TX tạo orđer

	discount_vl := 0.0
	total_amount := 0.0
	orderDetailTX := make([]services.OrderDetail, 0)
	orderDetailMOMO := make([]services.Product_MOMO, 0)
	order_id := uuid.New().String()
	// create and check orderDetailTX
	for _, sku := range product_sku {
		quantity := 0
		for _, od := range order.NumOfProducts {
			if sku.ProductSkuID == od.Product_sku_id {
				quantity = od.Amount
				break
			}
		}
		if quantity <= 0 {
			return nil, assets_services.NewError(400, fmt.Errorf("quantity of product: %s , product_sku: %s unvalid %d", sku.ProductsSpuID, sku.ProductSkuID, quantity))
		}
		if sku.SkuStock-int32(quantity) <= 0 {
			return nil, assets_services.NewError(400, fmt.Errorf("product: %s , product_sku: %s unstock ", sku.ProductsSpuID, sku.ProductSkuID))
		}
		orderDetailTX = append(orderDetailTX, services.OrderDetail{
			OrderDetailID: uuid.New().String(),
			Quantity:      int32(quantity),
			UnitPrice:     sku.Price,
			ProductSkuID:  sku.ProductSkuID,
			OrderID:       order_id,
		})
		if order.Payment_id == s.env.PaymentOnline {
			orderDetailMOMO = append(orderDetailMOMO, services.Product_MOMO{
				ID:          sku.ProductSkuID,
				Name:        sku.Name,
				Description: sku.ShortDescription,
				ImageURL:    s.env.PublicID + "/media/products?id=" + sku.Image,
				Price:       int64(sku.Price),
				Currency:    "VND",
				Quantity:    quantity,
				TotalPrice:  int64(sku.Price) * int64(quantity),
			})
		}
		total_amount += float64(quantity) * sku.Price
	}

	// kiem tra nguoi dung co dung vocher khong
	// // neu co dung => lay thong tin vorcher
	// // // kiem tra so luong neu con thoa mang thi thuc hien tru tien ra <else> thi tra ve loi discound kkhoong được sử dụng

	if order.Discount_Id != "" {
		dc, err := s.repository.Discount(ctx, order.Discount_Id)
		if err != nil {
			return nil, assets_services.NewError(400, fmt.Errorf("discount %s", err.Error()))
		}
		if dc.Amount <= 0 {
			return nil, assets_services.NewError(400, fmt.Errorf("discount: %s  unstock ", dc.DiscountCode))
		}
		if dc.EndDate.Before(time.Now()) {
			return nil, assets_services.NewError(400, fmt.Errorf("discount: %s  has started", dc.DiscountCode))
		}
		if dc.StartDate.After(time.Now()) {
			return nil, assets_services.NewError(400, fmt.Errorf("discount: %s  yet start", dc.DiscountCode))
		}
		// fmt.Printf("dc.dc.minOrderValue: %f, dc.total_amount: %f", dc.MinOrderValue, total_amount)
		if dc.MinOrderValue > total_amount {
			return nil, assets_services.NewError(400, fmt.Errorf("discount: %s  unvalid total_amount order must more than %f", dc.DiscountCode, dc.MinOrderValue))
		}
		discount_vl = dc.DiscountValue
	}
	paymentStatus := assets_services.OrderTable_PaymentStatus_ThanhToanTrucTiep
	orderStatus := assets_services.OrderTable_OrderStatus_ChoXacNhan
	if order.Payment_id == s.env.PaymentOnline {
		paymentStatus = assets_services.OrderTable_PaymentStatus_ChoThanhToan
	}
	// create orderTX
	orderTX := &services.Orders{
		OrderID:           order_id,
		TotalAmount:       total_amount - discount_vl,
		CustomerAddressID: order.Address_id,
		DiscountID:        services.Narg[string]{Data: order.Discount_Id, Valid: discount_vl != 0},
		PaymentMethodID:   order.Payment_id,
		CustomerID:        user_id,
		PaymentStatus:     paymentStatus,
		OrderStatus:       orderStatus,
	}
	err = s.repository.TXCreateOrdder(ctx, orderTX, orderDetailTX)
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("error when create order %s", err.Error()))
	}
	// sử lý thanh toán momo
	if order.Payment_id == s.env.PaymentOnline {
		payloadParamsMoMo := services.CombinedDataPayLoadMoMo{
			Info:    info,
			Address: address,
			OrderTX: orderTX,
			Items:   orderDetailMOMO,
		}
		payload := createMoMoPayload(s.env, payloadParamsMoMo)
		defer s.redis.AddOrderOnline(ctx, user_id, payloadParamsMoMo, s.env.OrderDuration)

		// log đầy đủ thông tin của payload dễ nhìn
		return callMoMoGetURL(s.env, payload)
	}

	return nil, nil
}
func (s *service) CallBackMoMo(ctx context.Context, tran services.TransactionMoMO) {
	// su ly api thanh cong
	// sử lý thoong báo tại đây
	if tran.ResultCode == 0 {
		orderID := tran.RequestID

		s.redis.DeleteOrderOnline(ctx, orderID)
		// s.repository
		// caapj nhat cot
		s.repository.UpdateOrder(ctx, services.Orders{
			OrderID:       orderID,
			PaymentStatus: assets_services.OrderTable_PaymentStatus_DaThanhToan,
		})
		//get order
		order, _ := s.repository.GetOrderByID(ctx, orderID)

		// get infousser
		userInfo, _ := s.repository.GetCustomer(ctx, order.CustomerID)

		if userInfo.DeviceRegistrationToken.Data != "" {
			// send message to
			msg := services_assets_sendMessage.ThanhToanThanhCong(orderID, tran.Amount)
			// get INFO
			s.firebase.SendToToken(ctx, userInfo.DeviceRegistrationToken.Data, msg)

		}

	}
	// su ly api that bai
}
func (s *service) ListOrderByUserID(ctx context.Context, user_id string, query services.QueryFilter) (map[string]interface{}, *assets_services.ServiceError) {
	// query.Conditions = []services.Condition{
	// 	{Field: "amount", Operator: ">=", Value: 1},
	// 	{Field: "end_date", Operator: ">", Value: time.Now()},
	// }
	orders, totalPages, totalElements, err := s.repository.GetOrdersByUserID(ctx, user_id, query)
	if err != nil {
		fmt.Println("Error ListDiscount:", err)
		return nil, assets_services.NewError(400, err)
	}

	result, err := assets_services.HideFields(orders, "orders", "customer_id")
	if err != nil {
		fmt.Println("Error HideFields:", err)
		return nil, assets_services.NewError(400, err)
	}
	result["currentPage"] = query.Page
	result["totalPages"] = totalPages
	result["totalElements"] = totalElements
	result["limit"] = query.PageSize

	// result["orderDetail"] = orderDetail
	return result, nil
}

func callMoMoGetURL(env config_assets.ReadENV, payload services.Payload_MOMO) (map[string]interface{}, *assets_services.ServiceError) {

	var jsonPayload []byte
	var err error
	jsonPayload, err = json.Marshal(payload)
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("error when json.Marshal %s", err.Error()))
	}
	//send HTTP to momo endpoint
	resp, err := http.Post(env.EndPointMoMo, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("error when send HTTP to momo endpoint: %s", err.Error()))
	}

	//result
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// func có đơn hàng nào đang chưa được thanh toán hay không
// nếu có thì tạo url tới momo trả về để người dùng thanh toán
// nếu không có trả về 204
func (s *service) GetURLOrderMoMOAgain(ctx context.Context, user_id string) (map[string]interface{}, *assets_services.ServiceError) {
	// check neu co order thi khong cho nguoi dung thanh toan
	payloadParamsMoMo, err := s.redis.GetOrderOnline(ctx, user_id)
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("error redis.GetOrderOnline:  %s ", err.Error()))
	}
	payload := createMoMoPayload(s.env, *payloadParamsMoMo)
	return callMoMoGetURL(s.env, payload)
}

func (s *service) RemoveOrderOnline(ctx context.Context, orderIDs string) {

	s.repository.TXCancelOrder(ctx, orderIDs)

}

func (s *service) CancelOrder(ctx context.Context, user_id, orderID string) *assets_services.ServiceError {
	order, err := s.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		return assets_services.NewError(400, fmt.Errorf("error repository.GetOrderByID:  %s ", err.Error()))
	}
	if order.CustomerID != user_id {
		return assets_services.NewError(400, fmt.Errorf("bạn không có quyền hủy đơn hàng này"))
	}
	// chir  cho don hang co trang thai cho xac nhan thi moi cho huy, con lai khong cho huy
	if order.OrderStatus != assets_services.OrderTable_OrderStatus_ChoXacNhan {
		return assets_services.NewError(400, fmt.Errorf("bạn không thể thể hủy đơn hàng này, trạng thái không phù hợp "))
	}
	if order.PaymentMethodID == s.env.PaymentOnline {
		return assets_services.NewError(400, fmt.Errorf("bạn không thể thể hủy đơn hàng thanh toán online, vui lòng liên hệ nhân viên tư vấn đề sử lý"))
	}

	err = s.repository.TXCancelOrder(ctx, orderID)
	if err != nil {
		return assets_services.NewError(400, fmt.Errorf("lỗi trong lúc hủy đơn ", err.Error()))

	}
	return nil
}

///v1/user/info/create_customeraddress
///v1/user/info/create_customeraddress
