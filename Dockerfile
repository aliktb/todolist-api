FROM golang:1.25.8-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /todolist-api ./cmd/api



FROM alpine:3.23

RUN apk add --no-cache curl ca-certificates

WORKDIR /app

COPY --from=builder /todolist-api /usr/local/bin/todolist-api
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/docs /app/docs

CMD ["todolist-api"]
