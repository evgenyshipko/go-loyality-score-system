# CChat helper for business

Project for learning golang on Yandex-practicum .

### Get started


0. Start PostgreSQL database

```bash
docker run -d --name metrics-collector-pg -p 5433:5432 -e POSTGRES_PASSWORD=metrics -e POSTGRES_USER=metrics -e POSTGRES_DB=metrics postgres
```

1. Install [goose](https://github.com/pressly/goose?tab=readme-ov-file#up)

(For MacOs):
```bash
brew install goose
```

2. Run migrations

UP:
```bash
goose -dir internal/db/migrations postgres "postgres://metrics:metrics@localhost:5433/metrics?sslmode=disable" up 
```

Down:
```bash
goose -dir internal/db/migrations postgres "postgres://metrics:metrics@localhost:5433/metrics?sslmode=disable" down 
```

3. Run server

```bash
go run ./cmd/gophermart -d="postgres://metrics:metrics@localhost:5433/metrics?sslmode=disable" -m=false
```