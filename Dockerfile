# Builder
FROM golang:1.21.3-alpine AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./

# Build
RUN go build -ldflags="-w -s" -o floodcontrol.out ./cmd/app/main.go

# Runner
FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/floodcontrol.out .

# Run
CMD ["./floodcontrol.out"]
