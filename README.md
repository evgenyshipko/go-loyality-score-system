# Chat helper for business (based on [RAG](https://habr.com/ru/articles/779526/))

## [YOUTUBE PRESENTATION](https://youtu.be/WwrNv7uGQdQ)

Project for learning golang on Yandex-practicum, based on [technical requirements](https://docs.google.com/document/d/1jaTT-PtvjUUaQoy0HK6dxhJ1OUSWX7NcRoQiCHpzRPk/edit?usp=sharing).

### Get started

0. Start local PostgreSQL database in docker container

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

3. Create .env file in your project root directory. Environment example you can see in .env.sample

4. Run server

```bash
go run ./cmd
```