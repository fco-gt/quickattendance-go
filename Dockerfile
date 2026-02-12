FROM golang:1.25.6-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o worker ./cmd/worker/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/worker .

# Cambiamos permisos para el usuario no-root
RUN chown appuser:appgroup /app/server /app/worker
USER appuser

EXPOSE 8080

# Por defecto arranca el server, pero se puede sobreescrito en docker-compose
CMD ["./server"]
