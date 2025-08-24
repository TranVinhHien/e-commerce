# 🛒 E-Commerce Application

An end-to-end **E-Commerce platform** with a **Golang API backend**, **MySQL + Redis** for data management, and a **Flutter mobile app** for end-users.  
The application provides a seamless online shopping experience with product browsing, purchasing, payment integration, and AI-powered search.  

---

## 📹 Demo
 [Watch the Demo Video](https://youtu.be/EXCrOGR7VL4](https://drive.google.com/file/d/1-6y-Cfs3L3tQdjZzr435HjA_9I9fK0rR/view?usp=drive_link))

---

## 🚀 Features

- **Product Listing** – Browse available products with categories and filters.  
- **Product Details** – View detailed information including images, description, and price.  
- **Shopping & Checkout** – Add products to cart and place orders.  
- **Online Payment** – Integrated with **MoMo** payment gateway.  
- **Order Tracking** – Check order history and real-time status updates.  
- **Smart Search**  
  - 🔎 Search by **text** keywords.  
  - 📷 Search by **image** using a Deep Learning model.  
- **Secure Session & Caching** – Redis-based caching and session management.  

---

## 🛠 Tech Stack

- **Backend:** Golang (RESTful API)  
- **Database:** MySQL 8.3.0  
- **Cache & Session:** Redis  
- **Mobile App:** Flutter  
- **AI Model:** Deep Learning (image & text search)  
- **Payment Integration:** MoMo  

---

## ⚙️ Backend Installation (Golang API)

### Prerequisites

- **MySQL** (v8.3.0) and **Redis**  
- **Go** (v1.22.2)  

### Run MySQL & Redis with Docker

```bash
# MySQL
docker run --name mysql_c -e MYSQL_ROOT_PASSWORD=12345 -p 3306:3306 -d mysql:8.3.0

# Redis
docker run -d --name redis_c -p 6379:6379   -v /data/redis-data/:/data   -e REDIS_ARGS="--requirepass 12345 --appendonly yes"   redis:latest
```

> **Note:**  
> - If running MySQL locally, set password to `12345`, port `3306`, version `8.3.0`.  
> - For local Redis, use port `6379`.  

---

### Install Golang

```bash
wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz
nano ~/.bashrc
# Add the following at the end:
export PATH=$PATH:/usr/local/go/bin
```

---

### Initialize Project

```bash
go mod init
```

---

### Run the Backend

```bash
# Initialize database
make initdb

# Start MySQL container
make startdb

# Start Redis container
make startredis

# Create database tables
make createtb

# Run the API server
make run
```

---

## 📌 Notes

- If dependencies are missing, install them manually or use `go get`.  
- For Linux systems, you may need to prefix commands with `sudo`.  

---

## 📱 Mobile App

The mobile app is built using **Flutter**, providing a clean and responsive UI for shopping, payments, and order tracking.  
You can run it on both **Android** and **iOS** devices.  
