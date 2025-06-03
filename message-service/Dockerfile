FROM golang:1.22-alpine

# Установка зависимостей
RUN apk add --no-cache git

WORKDIR /app

# Кешируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Сборка бинарника
WORKDIR /app/message-service/cmd
RUN go build -o message-service

# Запуск
CMD ["./message-service"]