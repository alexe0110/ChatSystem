




###

```bash
docker compose up -d
goose -dir migrations postgres "postgresql://postgres:postgres@localhost:5432/chat_db?sslmode=disable" up
```