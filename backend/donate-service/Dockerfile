FROM golang:1.23-alpine


WORKDIR /app

# Зависимости
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Исходники
COPY . .

# Сборка бинарника
RUN go build -o donate-service ./cmd/main.go

# Копируем скрипт ожидания PostgreSQL
COPY wait-for-postgres.sh /wait-for-postgres.sh
RUN chmod +x /wait-for-postgres.sh

# Открытый порт (опционально)
EXPOSE 8082

# Старт через ожидание БД
CMD ["/wait-for-postgres.sh", "/app/donate-service"]

