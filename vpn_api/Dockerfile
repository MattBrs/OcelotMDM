# Stage 1: Build the Go application
FROM golang:1.24.6-alpine3.21 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o vpn-api .

FROM kylemanna/openvpn

RUN apk add --no-cache ca-certificates expect

COPY --from=builder /app/vpn-api /usr/local/bin/vpn-api
COPY ./entrypoint.sh /usr/local/bin/entrypoint.sh

# Copy startup script
RUN chmod +x /usr/local/bin/vpn-api
RUN chmod +x /usr/local/bin/entrypoint.sh

EXPOSE 8080
# CMD ["/usr/local/bin/vpn-api"]
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
