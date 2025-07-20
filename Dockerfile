FROM golang:1.24-alpine AS builder

WORKDIR /build
COPY . . 
RUN go mod download
RUN go build -o ./go-marketplace-service ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /build
COPY --from=builder /build/go-marketplace-service ./go-marketplace-service
COPY .env .env
CMD ["/build/go-marketplace-service"]

