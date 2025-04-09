package db

import (
	services "new-project/services/entity"
	"time"
)

func (u *Accounts) Convert() services.Accounts {
	return services.Accounts{
		AccountID:    u.AccountID,
		Username:     u.Username,
		Password:     u.Password,
		ActiveStatus: string(u.ActiveStatus.AccountsActiveStatus),
		CreateDate:   u.CreateDate.Time,
		UpdateDate:   services.Narg[time.Time]{Data: u.UpdateDate.Time, Valid: u.UpdateDate.Valid},
	}
}

// Category
func (u *Categorys) Convert() services.Categorys {
	return services.Categorys{
		CategoryID: u.CategoryID,
		Name:       u.Name,
		Key:        u.Key,
		Path:       u.Path,
		Parent:     services.Narg[string]{Data: u.Parent.String, Valid: u.Parent.Valid},
	}
}

// Customer
func (u *Customers) Convert() services.Customers {
	return services.Customers{
		CustomerID: u.CustomerID,
		Name:       u.Name,
		Email:      u.Email,
		Image:      services.Narg[string]{Data: u.Image.String, Valid: u.Image.Valid},
		Dob:        u.Dob.Time,
		Gender:     string(u.Gender.CustomersGender),
		AccountID:  u.AccountID,
		CreateDate: u.CreateDate.Time,
		UpdateDate: services.Narg[time.Time]{Data: u.UpdateDate.Time, Valid: u.UpdateDate.Valid},
	}
}

// DescriptionAttr
func (u *DescriptionAttr) Convert() services.DescriptionAttr {
	return services.DescriptionAttr{
		DescriptionAttrID: u.DescriptionAttrID,
		Name:              u.Name,
		Value:             u.Value,
		ProductsSpuID:     u.ProductsSpuID,
	}
}

// Discount
func (u *Discounts) Convert() services.Discounts {
	return services.Discounts{
		DiscountID:     u.DiscountID,
		DiscountCode:   u.DiscountCode,
		DiscountValue:  u.DiscountValue,
		StartDate:      u.StartDate,
		EndDate:        u.EndDate,
		MinOrderValue:  u.MinOrderValue.Float64,
		Amount:         u.Amount.Int32,
		StatusDiscount: string(u.StatusDiscount.DiscountsStatusDiscount),
		CreateDate:     u.CreateDate.Time,
		UpdateDate:     services.Narg[time.Time]{Data: u.UpdateDate.Time, Valid: u.UpdateDate.Valid},
	}
}

// Employee
func (u *Employees) Convert() services.Employees {
	return services.Employees{
		EmployeeID:  u.EmployeeID,
		Gender:      string(u.Gender.EmployeesGender),
		Dob:         u.Dob.Time,
		Name:        u.Name,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Address:     u.Address.String,
		Salary:      u.Salary.Float64,
		CreateDate:  u.CreateDate.Time,
		UpdateDate:  services.Narg[time.Time]{Data: u.UpdateDate.Time, Valid: u.UpdateDate.Valid},
		AccountID:   u.AccountID,
	}
}

// Order
func (u *Orders) Convert() services.Orders {
	return services.Orders{
		OrderID:         u.OrderID,
		OrderDate:       u.OrderDate.Time,
		TotalAmount:     u.TotalAmount,
		DiscountID:      services.Narg[string]{Data: u.DiscountID.String, Valid: u.DiscountID.Valid},
		PaymentMethodID: u.PaymentMethodID,
		PaymentStatus:   string(u.PaymentStatus.OrdersPaymentStatus),
		OrderStatus:     string(u.OrderStatus.OrdersOrderStatus),
		CreateDate:      u.CreateDate.Time,
		UpdateDate:      services.Narg[time.Time]{Data: u.UpdateDate.Time, Valid: u.UpdateDate.Valid},
		CustomerID:      u.CustomerID,
	}
}

// OrderDetail
func (u *OrderDetail) Convert() services.OrderDetail {
	return services.OrderDetail{
		OrderDetailID: u.OrderDetailID,
		Quantity:      u.Quantity,
		UnitPrice:     u.UnitPrice,
		ProductSkuID:  u.ProductSkuID,
		OrderID:       u.OrderID,
	}
}

// PaymentMethod
func (u *PaymentMethods) Convert() services.PaymentMethods {
	return services.PaymentMethods{
		PaymentMethodID: u.PaymentMethodID,
		Name:            u.Name,
		Description:     services.Narg[string]{Data: u.Description.String, Valid: u.Description.Valid},
	}
}

// ProductSku
func (u *ProductSkus) Convert() services.ProductSkus {
	return services.ProductSkus{
		ProductSkuID:  u.ProductSkuID,
		Value:         u.Value,
		SkuStock:      u.SkuStock.Int32,
		Price:         u.Price,
		Sort:          u.Sort.Int32,
		CreateDate:    u.CreateDate.Time,
		UpdateDate:    services.Narg[time.Time]{Data: u.UpdateDate.Time, Valid: u.UpdateDate.Valid},
		ProductsSpuID: u.ProductsSpuID,
	}
}

// ProductSkuAttr
func (u *ProductSkuAttrs) Convert() services.ProductSkuAttrs {
	return services.ProductSkuAttrs{
		ProductSkuAttrID: u.ProductSkuAttrID,
		Name:             u.Name,
		Value:            u.Value,
		Image:            services.Narg[string]{Data: u.Image.String, Valid: u.Image.Valid},
		ProductsSpuID:    u.ProductsSpuID,
	}
}

// ProductSpu
func (u *ProductsSpu) Convert() services.ProductsSpu {
	return services.ProductsSpu{
		ProductsSpuID:    u.ProductsSpuID,
		Name:             u.Name,
		Description:      u.Description,
		ShortDescription: u.ShortDescription,
		StockStatus:      string(u.StockStatus.ProductsSpuStockStatus),
		DeleteStatus:     string(u.DeleteStatus.ProductsSpuDeleteStatus),
		Sort:             u.Sort.Int32,
		CreateDate:       u.CreateDate.Time,
		UpdateDate:       services.Narg[time.Time]{Data: u.UpdateDate.Time, Valid: u.UpdateDate.Valid},
		Image:            u.Image,
		Media:            string(u.Media),
		Key:              u.Key,
		CategoryID:       u.CategoryID,
	}
}

// PurchaseOrder
func (u *PurchaseOrders) Convert() services.PurchaseOrders {
	return services.PurchaseOrders{
		PurchaseOrderID: u.PurchaseOrderID,
		TotalAmount:     u.TotalAmount,
		Status:          string(u.Status.PurchaseOrdersStatus),
		CreateDate:      u.CreateDate.Time,
		UpdateDate:      services.Narg[time.Time]{Data: u.UpdateDate.Time, Valid: u.UpdateDate.Valid},
		SupplierID:      u.SupplierID,
		EmployeeID:      u.EmployeeID,
	}
}

// PurchaseOrderDetail
func (u *PurchaseOrderDetail) Convert() services.PurchaseOrderDetail {
	return services.PurchaseOrderDetail{
		PurchaseOrderDetailID: u.PurchaseOrderDetailID,
		Quantity:              u.Quantity.Int32,
		UnitPrice:             u.UnitPrice.Float64,
		PurchaseOrderID:       u.PurchaseOrderID,
		ProductSkuID:          u.ProductSkuID,
	}
}

// Role
func (u *Roles) Convert() services.Roles {
	return services.Roles{
		RoleID:      u.RoleID,
		Name:        u.Name,
		Description: services.Narg[string]{Data: u.Description.String, Valid: u.Description.Valid},
	}
}

// RoleAccount
func (u *RoleAccount) Convert() services.RoleAccount {
	return services.RoleAccount{
		RoleAccountID: u.RoleAccountID,
		AccountID:     u.AccountID,
		RoleID:        u.RoleID,
	}
}

// Supplier
func (u *Suppliers) Convert() services.Suppliers {
	return services.Suppliers{
		SupplierID:  u.SupplierID,
		Name:        u.Name,
		PhoneNumber: u.PhoneNumber,
		Email:       u.Email,
		Address:     services.Narg[string]{Data: u.Address.String, Valid: u.Address.Valid},
		CreateDate:  u.CreateDate.Time,
		UpdateDate:  services.Narg[time.Time]{Data: u.UpdateDate.Time, Valid: u.UpdateDate.Valid},
	}
}
