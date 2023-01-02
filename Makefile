dev:
	PORT=2565 DATABASE_URL=postgresql://postgres:postgres@localhost/expenses?sslmode=disable re go run server.go