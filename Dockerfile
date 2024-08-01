# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /url-shortener main.go

# Stage 2: Copy the binary to a smaller image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /url-shortener .

EXPOSE 8080

CMD ["./url-shortener"]
