
# ==========================================
# BUILD
# ==========================================
FROM golang:1.27.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-s -w" \
    -o /app/server ./cmd

# ==========================================
# RUNTIME
# ==========================================
FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/server ./server

EXPOSE 8080

CMD ["./server"]
