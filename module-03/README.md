#ЗАДАНИЕ 

```dockerfile 
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o myapp app.go

FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup && apk add --no-cache curl

COPY --from=builder /app/myapp /app/myapp

RUN chown -R appuser:appgroup /app/myapp

USER appuser

HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
    cmd curl --fail http://localhost:5000/health || exit 1

EXPOSE 5000

CMD ["/app/myapp"]

```


##Задание с java

``` java 
FROM eclipse-temurin:21-jdk AS build
WORKDIR /app

COPY Main.java .

RUN javac Main.java

CMD ["java", "Main"]

 FROM eclipse-temurin:21-jre-alpine AS final

 WORKDIR /app

 COPY --from=build /app/*.class .


 CMD ["java", "Main"]

``` 