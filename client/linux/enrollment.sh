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

res="$(curl -H 'Content-Type: application/json' \
      --silent \
      -w "\n%{http_code}\n" \
      -d "{ \"otp\":\"$otp\",\"type\":\"$type\"}" \
      -X POST \
      http://localhost:8080/devices)"

# TODO: parse res by first reading the code and if it's 200
# save the following data into a configuration file
echo "res: $res"
