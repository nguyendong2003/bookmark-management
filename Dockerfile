# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bookmark_service cmd/api/main.go


# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bookmark_service .

EXPOSE 8080

CMD ["./bookmark_service"]