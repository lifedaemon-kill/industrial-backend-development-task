FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /app/main ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 50051 8080 8081

ENTRYPOINT ["./main"]
