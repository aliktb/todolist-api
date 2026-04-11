FROM golang:1.25.8-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/todolist-api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/todolist-migrate ./cmd/migrate

FROM alpine:3.23

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /out/todolist-api /usr/local/bin/todolist-api
COPY --from=builder /out/todolist-migrate /usr/local/bin/todolist-migrate
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/docs /app/docs
COPY docker/entrypoint.sh /usr/local/bin/docker-entrypoint.sh

RUN chmod +x /usr/local/bin/docker-entrypoint.sh

ENV APP_ADDR=:8080

EXPOSE 8080

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["todolist-api"]
