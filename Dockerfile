FROM golang:latest AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o auth cmd/auth/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth .
CMD ["./auth"]