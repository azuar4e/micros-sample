FROM golang:1.26.1-alpine AS builder
WORKDIR /src

COPY proto/go.mod proto/go.sum ./proto/
COPY product-service/go.mod product-service/go.sum ./product-service/
WORKDIR /src/product-service
RUN go mod download

WORKDIR /src
COPY proto/ ./proto/
COPY product-service/ ./product-service/

WORKDIR /src/product-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/main ./cmd/

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/main /main
EXPOSE 8081 8077
ENTRYPOINT ["/main"]