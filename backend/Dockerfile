FROM golang:1.24.6-alpine3.21 AS builder

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o backend ./cmd/server


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/backend .

EXPOSE 8080
CMD ["./backend"]
