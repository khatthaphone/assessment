
PORT=2565
DATABASE_URL=postgresql://postgres:postgres@127.0.0.1:5432/expenses?sslmode=disable

ENV=PORT=${PORT} DATABASE_URL=${DATABASE_URL}

dev:
	${ENV} re go run .

test:
	${ENV} go test -v -cover ./... -coverprofile coverage.out

test-integration:
	${ENV} go test -v -tags integration ./...

test-it-docker:
	docker compose -f docker-compose.test.yaml up -d db
	sleep 10
	docker compose -f docker-compose.test.yaml up it_tests
	docker compose -f docker-compose.test.yaml down