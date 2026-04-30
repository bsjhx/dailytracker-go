#!/bin/bash
# Helper script to create a new user for DailyTracker
# Usage: ./create-user.sh username password [host]

if [ "$#" -lt 2 ] || [ "$#" -gt 3 ]; then
    echo "Usage: ./create-user.sh username password [host]"
    echo ""
    echo "Arguments:"
    echo "  username  - Username for the new user"
    echo "  password  - Password for the new user"
    echo "  host      - Optional host (default: localhost:8080)"
    echo ""
    echo "Example:"
    echo "  ./create-user.sh admin mypassword"
    echo "  ./create-user.sh admin mypassword example.com:8080"
    exit 1
fi

username=$1
password=$2
host=${3:-localhost:8080}

echo "Creating user '$username' on $host..."
echo ""

# Call the API endpoint
response=$(curl -s -w "\n%{http_code}" -X POST "http://$host/api/users/create" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$username\",\"password\":\"$password\"}")

# Extract HTTP status code (last line)
http_code=$(echo "$response" | tail -n1)
# Extract response body (everything except last line)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ] || [ "$http_code" -eq 201 ]; then
    echo "✅ User created successfully!"
    echo ""
    echo "You can now login with:"
    echo "  Username: $username"
    echo "  Password: $password"
else
    echo "❌ Error creating user (HTTP $http_code)"
    echo ""
    echo "Response:"
    echo "$body"
    exit 1
fi
