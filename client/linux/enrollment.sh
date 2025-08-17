#! /bin/bash

# Check prerequisites
if ! command -v curl &> /dev/null; then
  echo "Error: curl is not installed."
  exit 1
fi

if ! command -v openvpn &> /dev/null; then
  echo "Error: OpenVPN client not found. Please install it first."
  exit 1
fi


# Greet user
cat ./art.txt

# Request login stuff
echo "Welcome to the client enrollment script for OcelotMDM"

echo "Please enter the OTP to authenticate"

read otp

if [ -z "$otp" ]; then
  echo "Please enter something..."
  exit 1
fi

echo "Please enter the type the device (ticket-burner, autist personal device, ecc)"

read type

if [ -z "$type" ]; then
    echo "Please enter something..."
  exit 1
fi

echo "Using $otp for enrollment of a device of type '$type'"

system_arch=$(uname -m)

res=$(curl -H 'Content-Type: application/json' \
      --silent \
      -w "\n%{http_code}" \
      -d "{ \"otp\":\"$otp\",\"type\":\"$type\",\"architecture\":\"$system_arch\"}" \
      -X POST -L -k \
      https://159.89.2.75/devices)


http_code=$(echo "$res" | tail -n1)
body=$(echo "$res" | head -n -1)
echo "HTTP status: $http_code"

if [ "$http_code" -eq 201 ]; then
    echo "Enrollment successful!"

    device_name=$(echo "$body" | sed -n 's/.*"name": *"\([^"]*\)".*/\1/p')
    echo "device name: $device_name"

    ovpn_file=$(echo "$body" | sed -n 's/.*"ovpn_file": *"\([^"]*\)".*/\1/p') 
    echo "$ovpn_file" | sed -e 's/\\u003c/</g' \
      -e 's/\\u003e/>/g' \
      -e 's#\\\/#/#g' \
      -e 's/\\"/"/g' \
      -e 's/\\r//g' \
      -e 's/\\n/\n/g' \
  > "${device_name}.ovpn"    

    echo "VPN configuration saved to ${device_name}.ovpn"
    echo "You can now connect using: openvpn ${device_name}.ovpn"

else
    echo "Enrollment failed with HTTP code: $http_code"
    echo "Response: $body"
    exit 1
fi



