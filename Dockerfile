# Build stage
FROM golang:1.23.6-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz | tar xvz


# Run stage
FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY go.mod .
COPY go.sum .
COPY db/migration ./migration
RUN go mod download
RUN chmod +x /app/wait-for.sh /app/start.sh


EXPOSE 8080
ENTRYPOINT ["/app/start.sh"]
CMD ["/app/main"]