#!/bin/bash

echo "Starting OpenVPN and VPN API services..."

echo "Starting OpenVPN server..."
ovpn_run &
OPENVPN_PID=$!

sleep 5

echo "Starting VPN API server..."
cd /app
./vpn-api &
API_PID=$!

shutdown() {
    echo "Shutting down services..."
    kill $API_PID 2>/dev/null
    kill $OPENVPN_PID 2>/dev/null
    exit 0
}

trap shutdown SIGTERM SIGINT

# Wait for either process to exit
wait $API_PID
wait $OPENVPN_PID
