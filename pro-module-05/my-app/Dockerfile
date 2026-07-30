# Этап 1: сборка бинарника
FROM golang:1.21 AS builder

WORKDIR /app

# Сначала копируем go.mod отдельно - для кэширования слоёв
COPY go.mod ./
RUN go mod download

COPY main.go ./

# Собираем статический бинарник (без зависимостей от С-библиотек)
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Этап 2: минимальный финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем скомпилированный бинарник из этапа сборки
COPY --from=builder /app/server .

# Копируем статику (index.html и т.д.) - её тоже нужно явно скопировать
COPY static ./static

EXPOSE 8080

CMD ["./server"]