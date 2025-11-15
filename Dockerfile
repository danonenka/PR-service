# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Устанавливаем GOTOOLCHAIN для поддержки более новых версий
ENV GOTOOLCHAIN=auto

# Копируем go mod файлы
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем бинарник из builder stage
COPY --from=builder /app/main .

# Копируем openapi.yaml
COPY openapi.yaml .

EXPOSE 8080

CMD ["./main"]

