# ADR-002: Repository pattern с переключаемыми реализациями

## Status

Accepted

## Context

Приложение должно работать в двух окружениях: локально (docker-compose
с PostgreSQL и MinIO) и в AWS (DynamoDB и S3). Нужен способ переключения
между реализациями без изменения бизнес-логики.

Варианты:
- Build tags (`//go:build aws`) — разные сборки для разных окружений
- Одна сборка, runtime-переключение через env-переменные
- Абстрактная фабрика / DI-контейнер

## Decision

Runtime-переключение через env-переменную `DB_TYPE`. Интерфейсы определены
в `internal/repository/` (`UserRepository`, `MessageRepository`), реализации —
в подпакетах `postgres/` и `dynamodb/`.

Аналогично для хранилища файлов: интерфейс `FileStorage` в `pkg/storage/`,
реализации `MinioStorage` и `S3Storage`, переключение через `STORAGE_TYPE`.

DI ручная — через конструкторы в `main.go`.

## Consequences

**Плюсы:**
- Одна сборка для всех окружений — проще CI/CD
- Локальная разработка полностью автономна (docker-compose, без AWS)
- Добавление новой реализации (например, SQLite для тестов) — один пакет + одна ветка в main.go
- Service и handler слои не знают о конкретной БД

**Минусы:**
- Сборка содержит обе SDK (aws-sdk-go-v2 + lib/pq) — больше размер
- Handler зависит от конкретного `*service.UserService`, а не от интерфейса — нарушает симметрию pattern-а, затрудняет мокирование