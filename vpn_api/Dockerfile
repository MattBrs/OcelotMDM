# Stage 1: Build the Go application
FROM golang:1.24.6-alpine3.21 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vpn-api .

FROM kylemanna/openvpn

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/vpn-api /usr/local/bin/vpn-api

# Copy startup script
COPY start_services.sh /usr/local/bin/start_services.sh
RUN chmod +x /usr/local/bin/start_services.sh /usr/local/bin/vpn-api

EXPOSE 8080

CMD ["/usr/local/bin/start_services.sh"]
