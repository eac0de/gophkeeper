# Указываем базовый образ с поддержкой Go
FROM golang:1.23.3 AS builder

# Устанавливаем рабочую директорию в контейнере
WORKDIR /auth

# Копируем все файлы проекта в контейнер
COPY . .

# Загружаем зависимости (go.mod и go.sum) и устанавливаем их
RUN go mod download

# Собираем бинарный файл
RUN go build -o auth cmd/auth/main.go

# Используем минимальный образ для запуска собранного приложения
FROM debian:bookworm-slim

# Устанавливаем рабочую директорию для запуска приложения
WORKDIR /auth

# Копируем бинарный файл из стадии сборки
COPY --from=builder /auth/auth .
COPY --from=builder /auth/migrations /auth/migrations




