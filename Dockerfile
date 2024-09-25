FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o watchlist-app ./cmd/watchlist

FROM alpine:latest
RUN apk --no-cache add ca-certificates


WORKDIR /root/

COPY --from=builder /app/watchlist-app .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/static ./static

RUN mkdir certs


CMD ["./watchlist-app"]
