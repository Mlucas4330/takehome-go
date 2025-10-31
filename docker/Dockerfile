FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache git build-base

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

ENV CGO_ENABLED=0

CMD ["air"]
