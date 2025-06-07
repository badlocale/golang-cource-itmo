FROM golang:1.24.3-alpine AS builder

WORKDIR /app
COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o http-server ./cmd/http

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/http-server .
EXPOSE 8080

CMD ["./http-server"]