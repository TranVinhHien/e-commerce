Set up new project golang api 
yêu cầu cài đặt :
    - Đã có môi trường của mysql và redis .(có thể cài trên máy hoặc cài trên docker)
    lệnh để cài 2 môi trưởng trên trên docker:
        -mysql:docker run --name mysql_c -e MYSQL_ROOT_PASSWORD=12345 -p 3306:3306 -d mysql:8.3.0
    ( Nếu chạy mysql trên máy gốc thì yêu cầu phải đặt mật khẩu là 12345, port 3306, phiên bảng 8.3.0)
        -redis:docker run -d --name redis_c -p 6379:6379 -v /data/redis-data/:/data -e REDIS_ARGS="--requirepass 12345 --appendonly yes" redis:latest
    (Nếu chạy redis local thì phải đặt ở port 6379)
    - Cài đặt go trên môi trường máy tính có thể là linux hoặc windown , tải go phiên bản :go 1.22.2
        Lệnh tải trên linux:
        wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz
        sudo tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz
        nano ~/.bashrc
        thêm dòng sau vào cuối file:
        export PATH=$PATH:/usr/local/go/bin
    - Khi đã cài đặt môi trường go thực hiện lệnh : go mod init
    Khởi chạy ứng dụng như sau:  ( nếu làm trên linux thì thêm sudo ở trước)
        -Khởi tạo database: make initdb 
        -Khởi chạy container mysql: make startdb
        -Khởi chạy container redis: make startredis
        -Khởi tạo Bảng trong database: make createtb
        -Khởi chạy ứng dụng go bằng lệnh : make run
    Lưu ý trong quá trình chạy nếu gặp lỗi thiếu thư viện thì hãy tải hoặc nhờ chatgpt hỗ trợ.
