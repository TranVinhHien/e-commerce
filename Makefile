run:
	go run main.go
createredis:
	docker run -d --name redis_c -p 6379:6379 -v /data/redis-data/:/data -e REDIS_ARGS="--requirepass 12345 --appendonly yes" redis:latest
dropredis:
	docker rm redis_c
startredis:
	docker start redis_c
stopredis:
	docker stop redis_c
sqlc:
	sqlc generate
initmg:
	migrate create -ext sql -dir db/migration/ -seq init_mg
createc:
	docker run --name mysql_c -e MYSQL_ROOT_PASSWORD=12345 -p 3306:3306 -d mysql:8.3.0
rmc:
	docker rm mysql_c
initdb:
	docker exec -it mysql_c mysql -u root -p12345 -e "CREATE DATABASE \`e-commerce\`;"
dropdb:
	docker exec -it mysql_c mysql -u root -p12345 -e "DROP DATABASE \`e-commerce\`;"
startdb:
	docker start mysql_c
stopdb:
	docker stop mysql_c
buildimg:
	docker build -t tranvinhhien1912/e-commerce:$(tag) 
pushimg:
	docker push tranvinhhien1912/e-commerce:$(tag)
createtb:
	migrate -path db/migration/ -database "mysql://root:12345@tcp(localhost:3306)/e-commerce" -verbose up
droptb:
	migrate -path db/migration/ -database "mysql://root:12345@tcp(localhost:3306)/e-commerce" -verbose down
.PHONY: run