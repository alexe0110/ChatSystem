# ADR-001: Два API — REST + gRPC из общей бизнес-логики

## Status

Accepted

## Context

Чат-система требует два типа взаимодействия: классический CRUD (регистрация,
логин, отправка сообщений) и real-time коммуникацию (отправка сообщений,
загрузка файлов чанками и тд). HTTP/REST хорош для CRUD, для стриминга лучше
подходит gRPC.

Варианты:
- Только REST + WebSocket для real-time
- Только gRPC для всего
- REST для CRUD + gRPC для стриминга

## Decision

Два отдельных сервера: REST (Gin, :8080) и gRPC (:50051). Оба используют
общий слой `internal/service/` и `internal/repository/`, различаются только
транспортным слоем (`internal/handler/rest/` и `internal/handler/grpc/`).

gRPC покрывает все 4 типа RPC в одном проекте:
- Unary: GetUser, SendMessage
- Server streaming: GetMessageHistory
- Client streaming: UploadFile
- Bidirectional: Chat (real-time через Hub)

## Consequences

**Плюсы:**
- REST — привычный API для фронтенда и curl-тестирования
- gRPC — нативные стримы без костылей, типизированный контракт через protobuf
- Общий service-слой исключает дублирование бизнес-логики
- Демонстрирует все 4 типа gRPC в одном проекте

**Минусы:**
- Два сервиса, два процесса, две точки аутентификации
- Небольшая асимметрия: Notification worker подключён только к REST, не к gRPC