FROM golang:1.27.1-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /VikaBot

# Копируем файлы модулей отдельно для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код проекта
COPY . .

# Собираем бинарник статический (CGO_ENABLED=0 полезно для alpine/scratch)
RUN CGO_ENABLED=0 GOOS=linux go build -o /VikaBot/bin/main .

# --- Этап 2: Финальный минимальный образ ---
FROM alpine:latest

WORKDIR /root/

# Устанавливаем корневые сертификаты (если приложение делает HTTPS-запросы)
RUN apk --no-cache add ca-certificates

# Копируем собранный бинарник из первого этапа
COPY --from=builder /VikaBot/bin/main .

# Открываем порт (если приложение — сервер)
EXPOSE 8080

# Команда запуска
CMD ["./main"]