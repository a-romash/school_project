# syntax=docker/dockerfile:1

FROM golang:1.22.1-alpine AS builder

# Set destination for COPY
WORKDIR /app

COPY . .

RUN go build -o bin/server ./cmd
RUN go mod download && go mod verify

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/bin/server .
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./server"]