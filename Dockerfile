FROM golang:1.19 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

COPY go.* .

RUN go mod download

# Копируем исходный код приложения
COPY backend backend

# Собираем бинарный файл приложения с отключением CGO
RUN CGO_ENABLED=0 go build -o server ./backend/

# Стадия запуска
FROM alpine

# Копируем бинарный файл из стадии сборки
COPY --from=builder /app/server /app/server

WORKDIR /app

# Задаем переменную окружения для порта HTTP
ENV HTTP_PORT 8080

EXPOSE 8080

# Запускаем приложение
ENTRYPOINT ["./server"]