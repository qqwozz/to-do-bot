FROM golang:1.21-alpine AS builder

WORKDIR /build/app
COPY app/go.mod app/go.sum* ./
RUN go mod download
COPY app/ .
RUN CGO_ENABLED=0 GOOS=linux go test -count=1 ./... && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bot ./cmd/bot/

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
RUN addgroup -S botgroup && adduser -S botuser -G botgroup
WORKDIR /app
COPY --from=builder /bot .
USER botuser
CMD ["./bot"]
