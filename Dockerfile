# Build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/server -v ./cmd/production

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/bin/server .
COPY letters letters
ENTRYPOINT ["./server"]

LABEL Name=btcapp Version=0.0.1
EXPOSE 8080