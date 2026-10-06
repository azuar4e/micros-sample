FROM golang:1.26.1-alpine AS builder
WORKDIR /src

COPY proto/go.mod proto/go.sum ./proto/
COPY order-service/go.mod order-service/go.sum ./order-service/
WORKDIR /src/order-service
RUN go mod download

WORKDIR /src
COPY proto/ ./proto/
COPY order-service/ ./order-service/

WORKDIR /src/order-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/main ./cmd/

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/main /main
EXPOSE 8082
ENTRYPOINT ["/main"]