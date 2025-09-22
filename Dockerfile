FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/api/main.go

# Development
FROM golang:1.25.1-alpine AS development

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/air-verse/air@latest

COPY . .

EXPOSE 8080
CMD ["air"]

# Production
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/internal/infra/database/migrations ./internal/infra/database/migrations

EXPOSE 8080
CMD ["./main"]