# ---------- Builder ----------
FROM golang:1.25.6-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -ldflags="-w -s" -o server ./cmd/server
RUN go build -ldflags="-w -s" -o worker ./cmd/worker


# ---------- Runtime Final ----------
FROM alpine:3.20

RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/worker .

RUN chown appuser:appgroup /app/server /app/worker

USER appuser

EXPOSE 8080

# Por defecto arranca el server, Railway sobrescribirá esto para el Worker
ENTRYPOINT ["./server"]
