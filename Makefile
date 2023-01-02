
PORT=2565
DB_URL=postgresql://postgres:postgres@localhost/expenses?sslmode=disable
ENV=PORT=${PORT} DATABASE_URL=${DB_URL}

dev:
	${ENV} re go run server.go

test:
	${ENV} go test -v ./...

test-integration:
	${ENV} go test -v -tags integration ./...