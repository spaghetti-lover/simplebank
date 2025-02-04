postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine
createdb:
	docker exec -it postgres createdb --username=root --owner=root bank-system

dropdb:
	docker exec -it postgres dropdb bank-system

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank-system?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank-system?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank-system?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank-system?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate
test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/spaghetti-lover/bank-system/db/sqlc Store
	
.PHONY: postgres migrateup migrateup1 migratedown migratedown1 createdb dropdb sqlc test server mock