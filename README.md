# ChatSystem

Real-time чат на Go с деплоем в AWS через Terraform

Демонстрирует REST API, все 4 типа gRPC, конкурентность (goroutines/channels), Clean Architecture и cloud-native инфраструктуру.

## Архитектура

```mermaid
flowchart LR
    Client([Клиент])

    subgraph AWS ["AWS eu-central-1"]
        subgraph VPC ["VPC 10.0.0.0/16"]
            ALB[ALB :80]
            ECS[ECS Fargate :8080]
            ALB --> ECS
        end

        ECR[(ECR)]
        DynamoDB[(DynamoDB)]
        S3[(S3)]
        CW[CloudWatch]

        ECR -.->|pull| ECS
        ECS --> DynamoDB
        ECS --> S3
        ECS -.->|logs| CW
    end

    Client -->|HTTP| ALB
    GHA[GitHub Actions] -->|OIDC + push| ECR
```

## Технологии

| Слой        | Технология                                               |
|-------------|----------------------------------------------------------|
| Язык        | Go                                                       |
| REST        | Gin                                                      |
| gRPC        | Unary, Server streaming, Client streaming, Bidirectional |
| БД          | DynamoDB (AWS) / PostgreSQL (для локального запуска)     |
| Файлы       | S3 (AWS) / MinIO (локально)                              |
| Авторизация | JWT (REST middleware + gRPC interceptors)                |
| IaC         | Terraform                                                |
| CI/CD       | GitHub Actions                                           |
| Логи        | CloudWatch                                               |

## Структура проекта

```
chat-system/
├── cmd/
│   ├── rest/main.go                  # REST-сервер
│   └── grpc/main.go                  # gRPC-сервер
├── internal/
│   ├── handler/rest/                 # Gin handlers
│   ├── handler/grpc/                 # gRPC implementations
│   ├── service/                      # Бизнес-логика
│   ├── repository/                   # Интерфейсы доступа к данным
│   │   ├── postgres/                 # PostgreSQL-реализация
│   │   └── dynamodb/                 # DynamoDB-реализация
│   ├── hub/                          # Маршрутизация real-time сообщений
│   ├── middleware/                    # JWT auth (REST + gRPC)
│   ├── model/                        # Доменные модели
│   └── worker/                       # Фоновый notification worker
├── pkg/
│   ├── auth/                         # JWT утилиты
│   └── storage/                      # FileStorage интерфейс + S3/MinIO
├── terraform/                        # Инфраструктура
│   ├── main.tf                       # Provider, ECR, DynamoDB, S3
│   ├── vpc.tf                        # VPC, подсети, internet gateway
│   ├── sg.tf                         # Security groups
│   ├── alb.tf                        # Application Load Balancer
│   ├── ecs.tf                        # Fargate task и service
│   ├── iam.tf                        # IAM роли, политики, OIDC
│   ├── variables.tf                  # Входные переменные
│   └── outputs.tf                    # Выходные значения
├── .github/workflows/deploy.yml      # CI/CD pipeline
├── proto/                            # Protobuf файлы
├── migrations/                       # PostgreSQL миграции
├── docs/adr/                         # Architecture Decision Records
├── Dockerfile                        # Сборка образов
└── compose.yml                       # Для локальной разработки
```

## Переменные окружения

Приложение переключается между AWS и локальным режимом через env-переменные, 
- для локальной разработки нужно определить файл .env
- для AWS они определены в terraform/ecs.tf

| Переменная       | Описание                | AWS            | Локально                                                                |
|------------------|-------------------------|----------------|-------------------------------------------------------------------------|
| `DB_TYPE`        | Тип базы данных         | `dynamodb`     | -                                                                       |
| `STORAGE_TYPE`   | Тип файлового хранилища | `s3`           | -                                                                       |
| `AWS_REGION`     | Регион AWS              | `eu-central-1` | -                                                                       |
| `DATABASE_URL`   | PostgreSQL DSN          | -              | `postgresql://postgres:postgres@localhost:5432/chat_db?sslmode=disable` |
| `MINIO_ENDPOINT` | MinIO адрес             | -              | `localhost:9000`                                                        |
| `MINIO_USER`     | MinIO пользователь      | -              | `minioadmin`                                                            |
| `MINIO_PASSWORD` | MinIO пароль            | -              | `minioadmin`                                                            |

## Инфраструктура

Вся инфраструктура описана в Terraform

| Сервис      | Назначение                                 |
|-------------|--------------------------------------------|
| ECS Fargate | Запуск контейнера без управления серверами |
| ALB         | HTTP-балансировщик с health checks         |
| DynamoDB    | NoSQL БД, on-demand billing                |
| S3          | Файловое хранилище                         |
| ECR         | Приватный Docker registry                  |
| IAM         | Роли с минимальными правами                |
| CloudWatch  | Централизованные логи                      |
| VPC         | Изолированная сеть, 2 AZ                   |

**Стоимость:**  Балансировщик ALB не бесплатный, там $0.03/ч, чтобы минимизировать расходы `terraform apply` / `terraform destroy`

## Деплой

### В AWS

```bash
# 1. Поднять инфру в AWS
cd terraform
terraform apply

# 2. Получить ECR_URL
ECR_URL=$(terraform output -raw ecr_repository_url)
ACC_ID=${ECR_URL%%/*}

# 3. Собрать и запушить образ (без этого сервис будет 503)
aws ecr get-login-password --region eu-central-1 | docker login --username AWS --password-stdin $ACC_ID

docker build -t chat-system .
docker tag chat-system:latest ${ECR_URL}:latest
docker push ${ECR_URL}:latest

# 4. Получить URL
terraform output app_url

# <можно перейти по урлу и протестировать приложение>

# 5. Логи в реальном времени
aws logs tail /ecs/chatsystem --follow --region eu-central-1

# 6. Погасить ‼️ (чтобы минимум расходов)
terraform destroy
```

Если после пуша образа в ECR сервис все еще не стартует, то дело в backoff таймаутах на поиск образа.
Можно подождать или принудительно обновить сервис в ECS:
```bash
aws ecs update-service \
  --cluster chatsystem-cluster \
  --service chatsystem-service \
  --force-new-deployment \
  --region eu-central-1
```

### Локально

```bash
# Определить переменные в .env, затем определить их в окружении:
export $(grep -v '^#' .env | xargs)

# PostgreSQL + MinIO
docker compose up -d

# Миграции
goose -dir migrations postgres \
  "postgresql://postgres:postgres@localhost:5432/chat_db?sslmode=disable" up

# Запуск
go run cmd/rest/main.go
go run cmd/grpc/main.go
```

## API

### REST

```bash
# Регистрация
curl -X POST http://localhost:8080/user/register \
  -H "Content-Type: application/json" \
  -d '{"login":"alex","name":"Alex","password":"123456"}'

# Логин → JWT-токен
curl -X POST http://localhost:8080/user/login \
  -H "Content-Type: application/json" \
  -d '{"login":"alex","password":"123456"}'

# Получить юзера по ID
curl http://localhost:8080/user/<uuid> \
  -H "Authorization: Bearer <token>"

# Отправить сообщение
curl -X POST http://localhost:8080/message/send \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"sender_id":"<uuid>","receiver_id":"<uuid>","message_content":"Hello!"}'

# Получить сообщение по ID
curl http://localhost:8080/message/<message-uuid> \
  -H "Authorization: Bearer <token>"

# Получить переписку
curl "http://localhost:8080/message/conversation?sender_id=<uuid>&receiver_id=<uuid>" \
  -H "Authorization: Bearer <token>"
```

### gRPC (порт 50051)

**Unary - получить пользователя:**
```bash
grpcurl -plaintext -d '{
  "id": "185e2a36-da6c-44ec-b1ee-59d30bc3b948"
}' localhost:50051 chat.UserService/GetUser
```

**Unary - отправить сообщение:**
```bash
grpcurl -plaintext -d '{
  "sender_id": "185e2a36-da6c-44ec-b1ee-59d30bc3b948",
  "receiver_id": "9f4eba7f-96ab-4b72-b584-8853e4655010",
  "message_content": "Привет!"
}' localhost:50051 chat.ChatService/SendMessage
```

**Server streaming - история сообщений:**
```bash
grpcurl -plaintext -d '{
  "sender_id": "185e2a36-da6c-44ec-b1ee-59d30bc3b948",
  "receiver_id": "9f4eba7f-96ab-4b72-b584-8853e4655010"
}' localhost:50051 chat.ChatService/GetMessageHistory
```

**Client streaming - загрузка файла:**
```bash
# Отправляет чанки файла, получает ответ с URL (S3/MinIO)
grpcurl -plaintext -d @ localhost:50051 chat.ChatService/UploadFile <<EOF
{"file_name": "photo.jpg", "data": "<base64>"}
EOF
```

**Bidirectional streaming - real-time чат:**
```bash
grpcurl -plaintext -d @ localhost:50051 chat.ChatService/Chat
# Первое сообщение регистрирует в Hub:
# {"sender_id": "...", "receiver_id": "...", "content": "Привет!"}
```

## CI/CD

GitHub Actions при push в `main`:

1. Авторизация в AWS через **OIDC** (без хранения ключей)
2. Сборка Docker-образа
3. Push в ECR с тегами `latest` + commit SHA
4. *(Опционально, закомментировано)* Обновление ECS-сервиса

## Архитектурные решения

См. [docs/adr/](docs/adr/):

- [001 — Почему DynamoDB, а не PostgreSQL](docs/adr/001-why-dynamodb.md)
- [002 — Почему ECS Fargate, а не EKS](docs/adr/002-why-ecs-fargate.md)
- [003 — OIDC вместо Access Keys для CI/CD](docs/adr/003-oidc-over-access-keys.md)

## Линтеры

```bash
golangci-lint run --fix ./...
```