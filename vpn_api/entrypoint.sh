#!/usr/bin/env sh
set -e

# Start the OpenVPN server using the helper from kylemanna/openvpn
# Requires /etc/openvpn to be initialized (ovpn_genconfig, ovpn_initpki) and mounted.
ovpn_run &

# Start the certificate API
exec /usr/local/bin/vpn-api
