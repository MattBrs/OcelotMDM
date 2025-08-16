FROM kylemanna/openvpn

# Install Go
RUN apk add --no-cache go git

# Set Go environment
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

# Create app directory
WORKDIR /app

# Copy Go modules first (for better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the Go application
RUN go build -o vpn_api .

# Copy startup script
COPY start_services.sh /usr/local/bin/start_services.sh
RUN chmod +x /usr/local/bin/start_services.sh

EXPOSE 8080

CMD ["/usr/local/bin/start_services.sh"]
