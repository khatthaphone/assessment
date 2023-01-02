
PORT=2565
DB_URL=postgresql://postgres:postgres@localhost/expenses?sslmode=disable

dev:
	PORT=${PORT} DATABASE_URL=${DB_URL} re go run server.go

test:
	PORT=${PORT} DATABASE_URL=${DB_URL} go test -v ./...

test-integration:
	PORT=${PORT} DATABASE_URL=${DB_URL} go test -v -tags integration ./...