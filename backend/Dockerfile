# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o backend-server ./main.go

# Stage 2: Create a minimal image to run the binary
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/backend-server .
EXPOSE 8080
CMD ["./backend-server"]

