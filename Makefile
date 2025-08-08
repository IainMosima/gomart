postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=supersecret -d postgres:14-alpine

createdb:
	psql postgres://root:supersecret@localhost:5432/postgres -c "CREATE DATABASE gomart_db;"

dropdb:
	psql postgres://root:supersecret@localhost:5432/postgres -c "DROP DATABASE IF EXISTS gomart_db;"

migrateup:
	migrate -path infrastructures/db/migration -database postgres://root:supersecret@localhost:5432/gomart_db?sslmode=disable -verbose up

migratedown:
	migrate -path infrastructures/db/migration -database postgres://root:supersecret@localhost:5432/gomart_db?sslmode=disable -verbose down


.PHONY: postgres createdb dropdb migratedown migratedown