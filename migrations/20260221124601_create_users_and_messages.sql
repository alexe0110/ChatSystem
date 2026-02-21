-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE Users (
    id UUID PRIMARY KEY,
    login Text not null,
    name Text,
    hashed_password Text not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) IF NOT EXISTS;

CREATE TABLE Messages (
    id UUID PRIMARY KEY,
    sender_id UUID,
    receiver_id UUID,
    message_content Text,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX sender_index on Messages (sender_id);
CREATE INDEX receiver_index on Messages (receiver_id);
CREATE INDEX sender_receiver_index on Messages (sender_id, receiver_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Users;
DROP TABLE Messages;
DROP INDEX sender_index;
DROP INDEX receiver_index;
-- +goose StatementEnd
