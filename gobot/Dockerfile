FROM golang:1.26-alpine AS builder

WORKDIR /app
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
# COPY .env.example .env

EXPOSE 8000

CMD ["main"]