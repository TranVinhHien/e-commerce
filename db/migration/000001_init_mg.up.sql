-- Thiết lập mã hóa UTF-8
SET NAMES utf8mb4;

-- 1. Bảng category 
CREATE TABLE categorys (
  category_id VARCHAR(36) PRIMARY KEY,
  name NVARCHAR(128) NOT NULL,
  `key` VARCHAR(128) NOT NULL,
  `path` VARCHAR(128) NOT NULL,
  `parent` VARCHAR(36),
  FOREIGN KEY (parent) REFERENCES categorys(category_id)
);

-- 2. Bảng supplier 
CREATE TABLE suppliers (
  supplier_id VARCHAR(36) PRIMARY KEY,
  name NVARCHAR(128) NOT NULL,
  phone_number VARCHAR(15) NOT NULL UNIQUE,
  email VARCHAR(254) NOT NULL UNIQUE,
  address NVARCHAR(500),
  create_date DATETIME DEFAULT NOW(),
  update_date DATETIME DEFAULT NOW()
);


-- 4. Bảng account 
CREATE TABLE accounts (
  account_id VARCHAR(36) PRIMARY KEY,
  username VARCHAR(50) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  active_status ENUM('Active', 'Inactive') DEFAULT 'Active',
  create_date DATETIME DEFAULT NOW(),
  update_date DATETIME DEFAULT NOW()
);

-- 3. Bảng customer 
CREATE TABLE customers (
  customer_id VARCHAR(36) PRIMARY KEY,
  name NVARCHAR(128) NOT NULL,
  email VARCHAR(254) NOT NULL UNIQUE,
  image VARCHAR(500) ,
  dob DATE,
  gender ENUM('Nam', 'Nữ'),
  account_id VARCHAR(36) NOT NULL,
  create_date DATETIME DEFAULT NOW(),
  update_date DATETIME DEFAULT NOW(),
  FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

-- 3.1. Bảng customer
CREATE TABLE customer_address (
  id_address VARCHAR(36) PRIMARY KEY,
  customer_id VARCHAR(36) NOT NULL,
  address NVARCHAR(500) NOT NULL, 
  phone_number VARCHAR(15) NOT NULL,
  create_date DATETIME DEFAULT NOW(),
  update_date DATETIME DEFAULT NOW(),
  FOREIGN KEY (customer_id) REFERENCES customers(customer_id)
);

-- 5. Bảng employee
CREATE TABLE employees (
  employee_id VARCHAR(36) PRIMARY KEY,
  gender ENUM('Nam', 'Nữ'),
  dob DATE,
  name NVARCHAR(128) NOT NULL,
  email VARCHAR(254) NOT NULL UNIQUE,
  phone_number VARCHAR(15) NOT NULL UNIQUE,
  address VARCHAR(500),
  salary DOUBLE DEFAULT 0 CHECK (salary >= 0),
  create_date DATETIME DEFAULT NOW(),
  update_date DATETIME DEFAULT NOW(),
  account_id VARCHAR(36) NOT NULL,
  FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

-- 6. Bảng products_spu
CREATE TABLE products_spu (
  products_spu_id VARCHAR(36) PRIMARY KEY,
  name TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  brand_id VARCHAR(36) NOT NULL,
  description TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  short_description TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  stock_status ENUM('InStock', 'OutOfStock') DEFAULT 'InStock',
  delete_status ENUM('Active', 'Deleted') DEFAULT 'Active',
  sort INT DEFAULT 0,
  create_date DATETIME DEFAULT NOW(),
  update_date DATETIME DEFAULT NOW(),
  image VARCHAR(500) NOT NULL,
  media TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `key` VARCHAR(500) NOT NULL UNIQUE,
  category_id VARCHAR(36) NOT NULL,
  FOREIGN KEY (brand_id) REFERENCES suppliers(supplier_id),
  FOREIGN KEY (category_id) REFERENCES categorys(category_id)
);

-- 7. Bảng description_attr
CREATE TABLE description_attr (
  description_attr_id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  value NVARCHAR(500) NOT NULL,
  products_spu_id VARCHAR(36) NOT NULL,
  FOREIGN KEY (products_spu_id) REFERENCES products_spu(products_spu_id)
);

-- 8. Bảng product_sku
CREATE TABLE product_skus (
  product_sku_id VARCHAR(36) PRIMARY KEY,
  value VARCHAR(128) NOT NULL, 
  sku_stock INT DEFAULT 0 CHECK (sku_stock >= 0),
  price DOUBLE NOT NULL DEFAULT 0 CHECK (price >= 0),
  sort INT DEFAULT 0,
  create_date DATETIME DEFAULT NOW(),
  update_date DATETIME DEFAULT NOW(),
  products_spu_id VARCHAR(36) NOT NULL,
  FOREIGN KEY (products_spu_id) REFERENCES products_spu(products_spu_id)
);

CREATE TABLE ratings(
  rating_id VARCHAR(36) PRIMARY KEY,
  comment TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  star int NOT NULL CHECK (star >= 1 AND star <= 5),
  create_date DATETIME DEFAULT NOW(),
  update_date DATETIME DEFAULT NOW(),
  customer_id VARCHAR(36) NOT NULL,
  products_spu_id VARCHAR(36) NOT NULL,
  FOREIGN KEY (customer_id) REFERENCES customer(customer_id),
  FOREIGN KEY (products_spu_id) REFERENCES products_spu(products_spu_id)
);

-- 9. Bảng product_sku_attr
CREATE TABLE product_sku_attrs (
  product_sku_attr_id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  value NVARCHAR(500) NOT NULL,
  image VARCHAR(500),
  products_spu_id VARCHAR(36) NOT NULL,
  FOREIGN KEY (products_spu_id) REFERENCES products_spu(products_spu_id)
);

-- 10. Bảng payment_method 
CREATE TABLE payment_methods (
  payment_method_id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  description NVARCHAR(500)
);

-- 11. Bảng discount 
CREATE TABLE discounts (
  discount_id VARCHAR(36) PRIMARY KEY,
  discount_code VARCHAR(50) NOT NULL UNIQUE,
  discount_value DOUBLE NOT NULL DEFAULT 0 CHECK (discount_value >= 0),
  start_date DATETIME NOT NULL,
  end_date DATETIME NOT NULL,
  min_order_value DOUBLE DEFAULT 0 CHECK (min_order_value >= 0),
  amount INT DEFAULT 0 CHECK (amount >= 0),
  status_discount ENUM('Còn Hiệu Lực', 'Hết Hạn', 'Hết Lượt Sử Dụng') DEFAULT 'Còn Hiệu Lực',
  create_date DATETIME DEFAULT NOW(),
  update_date DATETIME DEFAULT NOW()
);

-- 12. Bảng purchase_order
CREATE TABLE purchase_orders (
  purchase_order_id VARCHAR(36) PRIMARY KEY,
  total_amount DOUBLE NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
  status ENUM('Đang Chờ', 'Đã Xác Nhận', 'Trì Hoãn', 'Đã Giao') DEFAULT 'Đang Chờ',
  create_date DATETIME DEFAULT NOW(),
  update_date DATETIME DEFAULT NOW(),
  supplier_id VARCHAR(36) NOT NULL,
  employee_id VARCHAR(36) NOT NULL,
  FOREIGN KEY (supplier_id) REFERENCES suppliers(supplier_id),
  FOREIGN KEY (employee_id) REFERENCES employees(employee_id)
);

-- 13. Bảng purchase_order_detail
CREATE TABLE purchase_order_detail (
  purchase_order_detail_id VARCHAR(36) PRIMARY KEY,
  quantity INT CHECK (quantity > 0),
  unit_price DOUBLE CHECK (unit_price >= 0),
  purchase_order_id VARCHAR(36) NOT NULL,
  product_sku_id VARCHAR(36) NOT NULL,
  FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(purchase_order_id),
  FOREIGN KEY (product_sku_id) REFERENCES product_skus(product_sku_id)
);

-- 14. Bảng order 
CREATE TABLE `orders` (
  order_id VARCHAR(36) PRIMARY KEY,
  order_date DATETIME DEFAULT NOW(),
  total_amount DOUBLE NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
  customer_address_id varchar(36) NOT NULL,
  discount_id VARCHAR(36),
  payment_method_id VARCHAR(36) NOT NULL,
  payment_status ENUM('Chờ Thanh Toán','Thanh Toán Trực Tiếp', 'Đã Thanh Toán') DEFAULT 'Chờ Thanh Toán',
  order_status ENUM('Đã Hủy', 'Chờ Xác Nhận', 'Đã Xác Nhận', 'Đã Giao Hàng') DEFAULT 'Chờ Xác Nhận',
  create_date DATETIME DEFAULT NOW(),
  update_date DATETIME DEFAULT NOW(),
  customer_id VARCHAR(36) NOT NULL,
  FOREIGN KEY (customer_address_id) REFERENCES customer_address(id_address),
  FOREIGN KEY (discount_id) REFERENCES discounts(discount_id),
  FOREIGN KEY (payment_method_id) REFERENCES payment_methods(payment_method_id),
  FOREIGN KEY (customer_id) REFERENCES customers(customer_id)
);

-- 15. Bảng order_detail (phụ thuộc order, product_sku)
CREATE TABLE order_detail (
  order_detail_id VARCHAR(36) PRIMARY KEY,
  quantity INT NOT NULL CHECK (quantity > 0),
  unit_price DOUBLE NOT NULL CHECK (unit_price >= 0),
  product_sku_id VARCHAR(36) NOT NULL,
  order_id VARCHAR(36) NOT NULL,
  FOREIGN KEY (product_sku_id) REFERENCES product_skus(product_sku_id),
  FOREIGN KEY (order_id) REFERENCES `orders`(order_id)
);

-- 16. Bảng role (độc lập)
CREATE TABLE roles (
  role_id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  description NVARCHAR(500)
);

-- 17. Bảng role_account (phụ thuộc account, role)
CREATE TABLE role_account (
  role_account_id VARCHAR(36) PRIMARY KEY,
  account_id VARCHAR(36) NOT NULL,
  role_id VARCHAR(36) NOT NULL,
  FOREIGN KEY (account_id) REFERENCES accounts(account_id),
  FOREIGN KEY (role_id) REFERENCES roles(role_id)
);

-- Chèn dữ liệu vào bảng category

-- Chèn dữ liệu vào bảng supplier

-- Chèn dữ liệu vào bảng customer

-- Chèn dữ liệu vào bảng account (phụ thuộc customer)
INSERT INTO accounts (account_id, username, password, active_status, create_date, update_date) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'nhanvien1', '$2b$12$wQczPcY/PorA6c3W1wTKtuDolIglUMix4KehgCnsf8PFh7SBP8R6i', 'Active', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440002', 'nhanvien2', '$2b$12$wQczPcY/PorA6c3W1wTKtuDolIglUMix4KehgCnsf8PFh7SBP8R6i', 'Active', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440003', 'nhanvien3', '$2b$12$wQczPcY/PorA6c3W1wTKtuDolIglUMix4KehgCnsf8PFh7SBP8R6i', 'Active', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440004', 'nhanvien4', '$2b$12$wQczPcY/PorA6c3W1wTKtuDolIglUMix4KehgCnsf8PFh7SBP8R6i', 'Active', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440005', 'nhanvien5', '$2b$12$wQczPcY/PorA6c3W1wTKtuDolIglUMix4KehgCnsf8PFh7SBP8R6i', 'Active', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440006', 'quanly1', '$2b$12$wQczPcY/PorA6c3W1wTKtuDolIglUMix4KehgCnsf8PFh7SBP8R6i', 'Active', NOW(), NOW());

-- Chèn dữ liệu vào bảng employee (phụ thuộc account)
INSERT INTO employees (employee_id, gender, dob, name, email, phone_number, address, salary, create_date, update_date, account_id) VALUES
('emp-003', 'Nam', '1996-07-10', 'Lê Văn C', 'nhanvien1@example.com', '0912345678', '789 Đường DEF, TP.HCM', 7500000, NOW(), NOW(), '550e8400-e29b-41d4-a716-446655440001'),
('emp-004', 'Nữ', '1998-12-25', 'Phạm Thị D', 'nhanvien2@example.com', '0913456789', '101 Đường GHI, TP.HCM', 7800000, NOW(), NOW(), '550e8400-e29b-41d4-a716-446655440002'),
('emp-005', 'Nam', '1994-03-15', 'Trương Văn E', 'nhanvien3@example.com', '0914567890', '202 Đường JKL, TP.HCM', 7600000, NOW(), NOW(), '550e8400-e29b-41d4-a716-446655440003'),
('emp-006', 'Nữ', '1997-09-30', 'Ngô Thị F', 'nhanvien4@example.com', '0915678901', '303 Đường MNO, TP.HCM', 7700000, NOW(), NOW(), '550e8400-e29b-41d4-a716-446655440004'),
('emp-007', 'Nam', '1993-11-05', 'Hoàng Văn G', 'nhanvien5@example.com', '0916789012', '404 Đường PQR, TP.HCM', 7900000, NOW(), NOW(), '550e8400-e29b-41d4-a716-446655440005'),
('emp-008', 'Nữ', '1988-06-20', 'Đỗ Thị H', 'quanly1@example.com', '0917890123', '505 Đường STU, TP.HCM', 16000000, NOW(), NOW(), '550e8400-e29b-41d4-a716-446655440006');
-- Chèn dữ liệu vào bảng payment_method
INSERT INTO payment_methods (payment_method_id, name, description) VALUES
('550e8400-e29b-41d4-a716-446655440050', 'Thanh toán online', 'Thanh toán qua thẻ tín dụng hoặc ví điện tử'),
('550e8400-e29b-41d4-a716-446655440051', 'Thanh toán khi nhận hàng', 'Thanh toán bằng tiền mặt khi nhận hàng');

-- Chèn dữ liệu vào bảng role
INSERT INTO roles (role_id, name, description) VALUES
('550e8400-e29b-41d4-a716-446655440060', 'Quản trị viên', 'Có quyền quản lý toàn bộ hệ thống'),
('550e8400-e29b-41d4-a716-446655440061', 'Nhân viên', 'Quản lý kiểm duyệt đơn hàng'),
('550e8400-e29b-41d4-a716-446655440062', 'Khách hàng', 'Người dùng thông thường');

-- Chèn dữ liệu vào bảng role_account (phụ thuộc account và role)
INSERT INTO role_account (role_account_id, account_id, role_id) VALUES
('ra-004', '550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440061'),
('ra-005', '550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440061'),
('ra-006', '550e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440061'),
('ra-007', '550e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440061'),
('ra-008', '550e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440061'),
('ra-009', '550e8400-e29b-41d4-a716-446655440006', '550e8400-e29b-41d4-a716-446655440060');