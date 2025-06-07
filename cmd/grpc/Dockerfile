FROM golang:1.24.3-alpine AS builder

WORKDIR /app
COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o grpc-server ./cmd/grpc

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/grpc-server .
EXPOSE 8081

CMD ["./grpc-server"]