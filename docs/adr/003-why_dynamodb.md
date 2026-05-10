# ADR-003: DynamoDB вместо managed PostgreSQL в AWS

## Status

Accepted

## Context

Для деплоя в AWS нужна БД:
- **RDS PostgreSQL** — знакомая реляционная модель, минимальные изменения в коде.
  Минимальная стоимость ~$15/мес (db.t3.micro).
- **DynamoDB** — NoSQL, on-demand billing, free tier 25 GB навсегда.
  Требует переосмысления модели данных.

Цель проекта — обучение

## Decision

DynamoDB с on-demand billing. Модель данных спроектирована от access patterns:

- `users`: PK = `user_id`, GSI `login-index` для поиска по логину при авторизации
- `messages`: PK = `chat_id` (отсортированная пара UUID), SK = `created_at`,
  GSI `message-id-index` для поиска по ID сообщения

`chat_id = sort(sender, receiver).join("#")` — обеспечивает симметричность:
диалог между A и B всегда имеет одинаковый ключ, независимо от того кто sender.

## Consequences

**Плюсы:**
- Бесплатно при низкой нагрузке
- Чтение по ключу — single-digit milliseconds, O(1)
- Не нужно управлять инстансом, бэкапами, патчами

**Минусы:**
- Нет JOIN-ов — денормализация обязательна
- Нет произвольных запросов без GSI (каждый новый access pattern = новый индекс)
- `GetMessageByID` требует отдельного GSI, тогда как в PostgreSQL это простой `WHERE id = $1`
- Нет foreign keys — целостность данных только на уровне приложения
