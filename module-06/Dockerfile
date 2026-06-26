FROM golang:1.26-alpine
WORKDIR /myapp
RUN go install github.com/air-verse/air@latest
COPY for_go/go.mod ./
RUN go mod download
CMD ["air"]