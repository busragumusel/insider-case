FROM golang:1.23.6 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o insider-case .

EXPOSE 8080
CMD ["./insider-case"]
