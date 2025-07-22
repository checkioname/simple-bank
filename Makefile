DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

network:
	docker network create bank-network

postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

test:
	go test -v -cover ./...

.PHONY: network postgres createdb dropdb migrateup migratedown test