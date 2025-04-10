package services

import (
	"time"
)

type Users struct {
}
type Accounts struct {
	AccountID    string          `json:"account_id"`
	Username     string          `json:"username"`
	Password     string          `json:"password"`
	ActiveStatus string          `json:"active_status"`
	CreateDate   time.Time       `json:"create_date"`
	UpdateDate   Narg[time.Time] `json:"update_date"`
}

type Categorys struct {
	CategoryID string            `json:"category_id"`
	Name       string            `json:"name"`
	Key        string            `json:"key"`
	Path       string            `json:"path"`
	Parent     Narg[string]      `json:"parent"`
	Childs     Narg[[]Categorys] `json:"child"`
}

type CustomerAddress struct {
	IDAddress   string          `json:"id_address"`
	CustomerID  string          `json:"customer_id"`
	Address     string          `json:"address"`
	PhoneNumber string          `json:"phone_number"`
	CreateDate  time.Time       `json:"create_date"`
	UpdateDate  Narg[time.Time] `json:"update_date"`
}

type Customers struct {
	CustomerID string          `json:"customer_id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Image      Narg[string]    `json:"image"`
	Dob        time.Time       `json:"dob"`
	Gender     string          `json:"gender"`
	AccountID  string          `json:"account_id"`
	CreateDate time.Time       `json:"create_date"`
	UpdateDate Narg[time.Time] `json:"update_date"`
}

type DescriptionAttr struct {
	DescriptionAttrID string `json:"description_attr_id"`
	Name              string `json:"name"`
	Value             string `json:"value"`
	ProductsSpuID     string `json:"products_spu_id"`
}

type Discounts struct {
	DiscountID     string          `json:"discount_id"`
	DiscountCode   string          `json:"discount_code"`
	DiscountValue  float64         `json:"discount_value"`
	StartDate      time.Time       `json:"start_date"`
	EndDate        time.Time       `json:"end_date"`
	MinOrderValue  float64         `json:"min_order_value"`
	Amount         int32           `json:"amount"`
	StatusDiscount string          `json:"status_discount"`
	CreateDate     time.Time       `json:"create_date"`
	UpdateDate     Narg[time.Time] `json:"update_date"`
}

type Employees struct {
	EmployeeID  string          `json:"employee_id"`
	Gender      string          `json:"gender"`
	Dob         time.Time       `json:"dob"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	PhoneNumber string          `json:"phone_number"`
	Address     string          `json:"address"`
	Salary      float64         `json:"salary"`
	CreateDate  time.Time       `json:"create_date"`
	UpdateDate  Narg[time.Time] `json:"update_date"`
	AccountID   string          `json:"account_id"`
}

type OrderDetail struct {
	OrderDetailID string  `json:"order_detail_id"`
	Quantity      int32   `json:"quantity"`
	UnitPrice     float64 `json:"unit_price"`
	ProductSkuID  string  `json:"product_sku_id"`
	OrderID       string  `json:"order_id"`
}

type Orders struct {
	OrderID           string          `json:"order_id"`
	OrderDate         time.Time       `json:"order_date"`
	TotalAmount       float64         `json:"total_amount"`
	CustomerAddressID string          `json:"customer_address_id"`
	DiscountID        Narg[string]    `json:"discount_id"`
	PaymentMethodID   string          `json:"payment_method_id"`
	PaymentStatus     string          `json:"payment_status"`
	OrderStatus       string          `json:"order_status"`
	CreateDate        time.Time       `json:"create_date"`
	UpdateDate        Narg[time.Time] `json:"update_date"`
	CustomerID        string          `json:"customer_id"`
}

type PaymentMethods struct {
	PaymentMethodID string       `json:"payment_method_id"`
	Name            string       `json:"name"`
	Description     Narg[string] `json:"description"`
}

type ProductSkuAttrs struct {
	ProductSkuAttrID string       `json:"product_sku_attr_id"`
	Name             string       `json:"name"`
	Value            string       `json:"value"`
	Image            Narg[string] `json:"image"`
	ProductsSpuID    string       `json:"products_spu_id"`
}

type ProductSkus struct {
	ProductSkuID  string          `json:"product_sku_id"`
	Value         string          `json:"value"`
	SkuStock      int32           `json:"sku_stock"`
	Price         float64         `json:"price"`
	Sort          int32           `json:"sort"`
	CreateDate    time.Time       `json:"create_date"`
	UpdateDate    Narg[time.Time] `json:"update_date"`
	ProductsSpuID string          `json:"products_spu_id"`
}

type ProductsSpu struct {
	ProductsSpuID    string          `json:"products_spu_id"`
	Name             string          `json:"name"`
	BrandID          string          `json:"brand_id"`
	Description      string          `json:"description"`
	ShortDescription string          `json:"short_description"`
	StockStatus      string          `json:"stock_status"`
	DeleteStatus     string          `json:"delete_status"`
	Sort             int32           `json:"sort"`
	CreateDate       time.Time       `json:"create_date"`
	UpdateDate       Narg[time.Time] `json:"update_date"`
	Image            string          `json:"image"`
	Media            string          `json:"media"`
	Key              string          `json:"key"`
	CategoryID       string          `json:"category_id"`
}

type PurchaseOrderDetail struct {
	PurchaseOrderDetailID string  `json:"purchase_order_detail_id"`
	Quantity              int32   `json:"quantity"`
	UnitPrice             float64 `json:"unit_price"`
	PurchaseOrderID       string  `json:"purchase_order_id"`
	ProductSkuID          string  `json:"product_sku_id"`
}

type PurchaseOrders struct {
	PurchaseOrderID string          `json:"purchase_order_id"`
	TotalAmount     float64         `json:"total_amount"`
	Status          string          `json:"status"`
	CreateDate      time.Time       `json:"create_date"`
	UpdateDate      Narg[time.Time] `json:"update_date"`
	SupplierID      string          `json:"supplier_id"`
	EmployeeID      string          `json:"employee_id"`
}

type Ratings struct {
	RatingID      string          `json:"rating_id"`
	Comment       Narg[string]    `json:"comment"`
	Star          int32           `json:"star"`
	CreateDate    time.Time       `json:"create_date"`
	UpdateDate    Narg[time.Time] `json:"update_date"`
	AccountID     string          `json:"account_id"`
	ProductsSpuID string          `json:"products_spu_id"`
}

type RoleAccount struct {
	RoleAccountID string `json:"role_account_id"`
	AccountID     string `json:"account_id"`
	RoleID        string `json:"role_id"`
}

type Roles struct {
	RoleID      string       `json:"role_id"`
	Name        string       `json:"name"`
	Description Narg[string] `json:"description"`
}

type Suppliers struct {
	SupplierID  string          `json:"supplier_id"`
	Name        string          `json:"name"`
	PhoneNumber string          `json:"phone_number"`
	Email       string          `json:"email"`
	Address     Narg[string]    `json:"address"`
	CreateDate  time.Time       `json:"create_date"`
	UpdateDate  Narg[time.Time] `json:"update_date"`
}
