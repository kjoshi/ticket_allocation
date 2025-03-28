DB_USER=postgres
DB_PASSWORD=postgres
API_DB_PORT=5432
API_DB_NAME=tickets
API_PORT=3000
TEST_DB_PORT=6543
TEST_DB_NAME=tickets-test

run:
	@DB_HOST=localhost DB_PORT=$(API_DB_PORT) DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) DB_NAME=$(API_DB_NAME) PORT=$(API_PORT) go run .

test:
	@DB_HOST=localhost DB_PORT=$(TEST_DB_PORT) DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) DB_NAME=$(TEST_DB_NAME) go test ./models

.PHONY: run test
