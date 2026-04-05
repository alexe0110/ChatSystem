




###

Подготовка окружения
```bash
docker compose up -d
goose -dir migrations postgres "postgresql://postgres:postgres@localhost:5432/chat_db?sslmode=disable" up
```

Объявить энвы из файла

```bash
export $(grep -v '^#' .env | xargs)
```

Запуск
```bash
go build ./cmd/rest/main.go
./main
---
go run main.go
```

Линтеры и форматеры
```bash
golangci-lint run --fix ./...
```


Создание новой миграции
```bash
goose -dir migrations create create_users_and_messages sql
```


### Тестирование

```bash
# 1. Регистрация
curl -X POST http://localhost:8080/user/register \
  -H "Content-Type: application/json" \
  -d '{"login":"alex","name":"Alex","password":"123456"}'

# 2. Логин (получить токен)
curl -X POST http://localhost:8080/user/login \
  -H "Content-Type: application/json" \
  -d '{"login":"alex","password":"123456"}'

# 3. Получить юзера по ID (подставь TOKEN и UUID)
curl http://localhost:8080/user/USER_UUID_HERE \
  -H "Authorization: Bearer TOKEN_HERE"

# 4. Отправить сообщение (нужен второй юзер)
curl -X POST http://localhost:8080/message/send \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN_HERE" \
  -d '{"sender_id":"SENDER_UUID","receiver_id":"RECEIVER_UUID","message_content":"Hello!"}'

# 5. Получить сообщение по ID
curl http://localhost:8080/message/MESSAGE_UUID_HERE \
  -H "Authorization: Bearer TOKEN_HERE"

# 6. Получить переписку
curl "http://localhost:8080/message/conversation?sender_id=SENDER_UUID&receiver_id=RECEIVER_UUID" \
  -H "Authorization: Bearer TOKEN_HERE"
```