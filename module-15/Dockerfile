FROM golang:1.26-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o myapp app.go

EXPOSE 5000
CMD ["./myapp"]