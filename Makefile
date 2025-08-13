postgres:
	docker run --name gomart-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=supersecret -d postgres:14-alpine

createdb:
	psql postgres://root:supersecret@localhost:5432/postgres -c "CREATE DATABASE gomart_db;"

dropdb:
	psql postgres://root:supersecret@localhost:5432/postgres -c "DROP DATABASE IF EXISTS gomart_db;"

migrateup:
	migrate -path infrastructures/db/migration -database postgres://root:supersecret@localhost:5432/gomart_db?sslmode=disable -verbose up

migratedown:
	migrate -path infrastructures/db/migration -database postgres://root:supersecret@localhost:5432/gomart_db?sslmode=disable -verbose down

sqlc:
	sqlc generate

mockgen:
	@echo "Usage: make mockgen SOURCE=path/to/interface.go DEST=mocks/mock_name.go"
	@echo "Example: make mockgen SOURCE=domains/category/repository/category_repo_int.go DEST=mocks/category_repo_mock.go"

test:
	go test ./... -v

test-coverage:
	go test ./... -cover

test-coverage-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-category:
	go test ./services/category/... ./infrastructures/repository/... ./rest-server/handlers/... -v

test-category-coverage:
	go test ./services/category/... ./infrastructures/repository/... ./rest-server/handlers/... -cover

clean-coverage:
	rm -f coverage.out coverage.html

seed-categories:
	psql postgres://root:supersecret@localhost:5432/gomart_db -f seed_categories.sql

seed-products:
	psql postgres://root:supersecret@localhost:5432/gomart_db -f seed_products.sql

seed: seed-categories seed-products
	@echo "Seed data loaded successfully!"

seed-verify:
	psql postgres://root:supersecret@localhost:5432/gomart_db -c "SELECT 'Categories' as type, COUNT(*)::text as count FROM categories WHERE is_deleted = FALSE UNION ALL SELECT 'Products' as type, COUNT(*)::text as count FROM products WHERE is_deleted = FALSE UNION ALL SELECT 'Price Range' as type, CONCAT('KES ', MIN(price), ' - ', MAX(price)) as count FROM products WHERE is_deleted = FALSE;"

.PHONY: postgres createdb dropdb migratedown migratedown sqlc mockgen test test-coverage test-coverage-html test-category test-category-coverage clean-coverage seed-categories seed-products seed seed-verify