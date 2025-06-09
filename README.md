# E-commerce Backend - Golang API

## Yêu cầu cài đặt

- **MySQL** và **Redis** (cài trên máy hoặc Docker)

### Cài đặt MySQL và Redis bằng Docker

```bash
# MySQL
docker run --name mysql_c -e MYSQL_ROOT_PASSWORD=12345 -p 3306:3306 -d mysql:8.3.0

# Redis
docker run -d --name redis_c -p 6379:6379 \
  -v /data/redis-data/:/data \
  -e REDIS_ARGS="--requirepass 12345 --appendonly yes" \
  redis:latest
```

> **Lưu ý:**  
> - Nếu chạy MySQL trên máy, đặt mật khẩu là `12345`, port `3306`, phiên bản `8.3.0`.  
> - Nếu chạy Redis local, sử dụng port `6379`.

---

### Cài đặt Golang

Tải và cài đặt Go phiên bản **1.22.2**:

```bash
wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz
nano ~/.bashrc
# Thêm dòng sau vào cuối file:
export PATH=$PATH:/usr/local/go/bin
```

---

## Khởi tạo dự án

Sau khi cài đặt Go, thực hiện lệnh:

```bash
go mod init
```

---

## Khởi chạy ứng dụng

> Nếu dùng Linux, thêm `sudo` trước các lệnh bên dưới nếu cần.

```bash
# Khởi tạo database
make initdb

# Khởi chạy container MySQL
make startdb

# Khởi chạy container Redis
make startredis

# Khởi tạo bảng trong database
make createtb

# Chạy ứng dụng Go
make run
```

---

## Lưu ý

- Nếu gặp lỗi thiếu thư viện, hãy cài đặt hoặc nhờ ChatGPT hỗ trợ. 