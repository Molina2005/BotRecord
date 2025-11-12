FROM golang:1.25 AS builder 
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 
RUN go build -o modulo .
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /app/modulo .
COPY .env .env
CMD ["./modulo"]
