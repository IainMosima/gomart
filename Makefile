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

mockgen-auth:
	mockgen -source=domains/auth/service/auth_service_int.go -destination=domains/auth/service/auth_service_mock.go -package=service
	mockgen -source=domains/auth/service/cognito_service_int.go -destination=domains/auth/service/cognito_service_mock.go -package=service

mockgen-category:
	mockgen -source=domains/category/service/category_service_int.go -destination=domains/category/service/category_service_mock.go -package=service

mockgen-product:
	mockgen -source=domains/product/service/product_service_int.go -destination=domains/product/service/product_service_mock.go -package=service

mockgen-customer:
	mockgen -source=domains/customer/service/customer_service_int.go -destination=domains/customer/service/customer_service_mock.go -package=service

mockgen-order:
	mockgen -source=domains/order/service/order_service_int.go -destination=domains/order/service/order_service_mock.go -package=service

mockgen-all: mockgen-auth mockgen-category mockgen-product mockgen-customer mockgen-order
	@echo "All domain mocks generated successfully!"

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

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

.PHONY: postgres createdb dropdb migratedown migratedown sqlc mockgen mockgen-auth mockgen-category mockgen-product mockgen-customer mockgen-order mockgen-all test test-coverage test-coverage-html test-category test-category-coverage clean-coverage seed-categories seed-products seed seed-verify docker-up docker-down docker-logs