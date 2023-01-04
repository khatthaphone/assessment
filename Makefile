
PORT=2565
DB_URL=postgresql://postgres:postgres@localhost/expenses?sslmode=disable
TEST_DB_URL=postgresql://postgres:postgres@db/expenses-it-test?sslmode=disable
ENV=PORT=${PORT} DATABASE_URL=${DB_URL}
TEST_ENV=PORT=${PORT} DATABASE_URL=${TEST_DB_URL}

dev:
	${ENV} re go run .

test:
	${TEST_ENV} go test -v ./...

test-integration:
	${ENV} go test -v -tags integration ./...

test-it-docker:
	docker compose -f docker-compose.test.yaml up -d db
	sleep 5
	docker compose -f docker-compose.test.yaml up it_tests
	docker compose -f docker-compose.test.yaml down --volumes