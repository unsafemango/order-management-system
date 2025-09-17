# syntax=docker/dockerfile:1
FROM golang:1.23.5
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/products

EXPOSE 8081
CMD ["./main"]