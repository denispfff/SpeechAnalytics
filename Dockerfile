FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web_server ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache tzdata && cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && echo "Europe/Moscow" > /etc/timezone

# Скопируем готовый бинарник из слоя сборки
COPY --from=builder /app/web_server .

# Открываем порт и задаем переменную окружения
EXPOSE 8080
ENV TODO_PORT="8080"

# Запускаем сервер с указанием команды help по умолчанию
ENTRYPOINT ["./web_server"]
CMD ["--help"]
