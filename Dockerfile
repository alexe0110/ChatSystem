FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o chat-server ./cmd/rest/


FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/chat-server .

CMD ["/app/chat-server"]