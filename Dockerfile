FROM golang:1.25.6-alpine

WORKDIR /opt/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bookmark_service cmd/api/main.go

CMD ["./bookmark_service"]