FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o mock-pipefy-api ./cmd/api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/mock-pipefy-api .

EXPOSE 8001

CMD ["./mock-pipefy-api"]