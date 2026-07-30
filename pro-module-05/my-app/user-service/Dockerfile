FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN adduser -D -g '' appuser
COPY --from=builder /app/server /usr/local/bin/server
RUN chmod +x /usr/local/bin/server
USER appuser
EXPOSE 8080
CMD ["/usr/local/bin/server"]