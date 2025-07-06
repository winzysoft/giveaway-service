FROM golang:1.24-alpine as builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY . .
RUN go mod download

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /main ./cmd/app

# Устанавливаем migrate в builder-стадии
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Final stage
FROM alpine:3.19

WORKDIR /app

# Копируем бинарники и миграции
COPY --from=builder /main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/config ./config
COPY --from=builder /go/bin/migrate /bin/migrate

# Устанавливаем только необходимые пакеты
RUN apk add --no-cache tzdata ca-certificates

# Настройки
ENV CONFIG_PATH=/app/config/local.yaml \
    TZ=UTC

EXPOSE 8080
