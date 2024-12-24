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
	docker run --name postgres_c -p 5432:5432  -e POSTGRES_USER=root -e POSTGRES_PASSWORD=12345 -d postgres:17-alpine
rmc:
	docker rm postgres_c
initdb:
	docker exec -it postgres_c createdb --username=root --owner=root new_project
dropdb:
	docker exec -it postgres_c dropdb new_project
startdb:
	docker start postgres_c
stopdb:
	docker stop postgres_c
buildimg:
	docker build -t new_project:latest .
pushimg:
	docker push tranvinhhien1912/new_project:tagname
createtb:
	migrate -path db/migration/ -database "postgresql://root:12345@localhost:5432/new_project?sslmode=disable" -verbose up
droptb:
	migrate -path db/migration/ -database "postgresql://root:12345@localhost:5432/new_project?sslmode=disable" -verbose down
.PHONY: run