# Stage 1: Build the Go application
FROM golang:1.24.6-alpine3.21 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vpn-api .

# Stage 2: Create the final image with OpenVPN
FROM kylemanna/openvpn

# Install any additional packages if needed
RUN apk add --no-cache ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/vpn-api /usr/local/bin/vpn-api

# Copy startup script
COPY start-services.sh /usr/local/bin/start-services.sh
RUN chmod +x /usr/local/bin/start-services.sh /usr/local/bin/vpn-api

EXPOSE 8080

CMD ["/usr/local/bin/start-services.sh"]
