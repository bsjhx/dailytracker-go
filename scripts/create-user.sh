#!/bin/bash
# Helper script to create a new user for DailyTracker
# Usage: ./create-user.sh username password

if [ "$#" -ne 2 ]; then
    echo "Usage: ./create-user.sh username password"
    exit 1
fi

username=$1
password=$2

# Generate bcrypt hash using Go
echo "Generating password hash..."

# Create a temporary Go file
tmpfile=$(mktemp /tmp/hash-password.XXXXXX.go)
cat > "$tmpfile" << 'EOF'
package main
import (
    "fmt"
    "os"
    "golang.org/x/crypto/bcrypt"
)
func main() {
    if len(os.Args) < 2 {
        fmt.Fprintf(os.Stderr, "Error: password argument required\n")
        os.Exit(1)
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(os.Args[1]), 10)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    fmt.Print(string(hash))
}
EOF

# Run the Go program with the password as an argument
hash=$(go run "$tmpfile" "$password" 2>&1)
exit_code=$?

# Clean up
rm "$tmpfile"

if [ $exit_code -ne 0 ]; then
    echo "Error generating password hash: $hash"
    exit 1
fi

echo ""
echo "Generated SQL to create user:"
echo "============================================"
echo "INSERT INTO users (username, password_hash)"
echo "VALUES ('$username', '$hash');"
echo "============================================"
echo ""
echo "To apply this, run:"
echo "  sqlite3 ./dailytracker.db"
echo ""
echo "Then paste the INSERT statement above."
