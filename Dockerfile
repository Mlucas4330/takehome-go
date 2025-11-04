FROM golang:1.25-alpine AS builder
WORKDIR /app

ENV CGO_ENABLED=0

RUN apk add --no-cache git build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /out/api ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=builder /out/api /app/api

EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/api"]