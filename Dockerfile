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


# ---------- Runtime base ----------
FROM alpine:3.20 AS runner-base

RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app
USER appuser


# ---------- Server ----------
FROM runner-base AS server
COPY --from=builder /app/server .
EXPOSE 8080
ENTRYPOINT ["./server"]


# ---------- Worker ----------
FROM runner-base AS worker
COPY --from=builder /app/worker .
ENTRYPOINT ["./worker"]
