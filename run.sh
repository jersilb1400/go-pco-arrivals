#!/bin/bash

# Load environment variables from .env if present
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

echo "🚀 Starting PCO Arrivals Billboard..."
echo "Using PCO_CLIENT_ID: $PCO_CLIENT_ID"

# Set default environment variables if not already set
export PCO_CLIENT_ID=${PCO_CLIENT_ID:-"test_client_id"}
export PCO_CLIENT_SECRET=${PCO_CLIENT_SECRET:-"test_client_secret"}
export PCO_REDIRECT_URI=${PCO_REDIRECT_URI:-"http://localhost:3000/auth/callback"}
export PORT=${PORT:-3000}
export HOST=${HOST:-"0.0.0.0"}

# Create data directory if it doesn't exist
mkdir -p data

# Build the application
echo "📦 Building application..."
go build -o pco-billboard .

# Run the application
echo "🌟 Starting server on http://localhost:$PORT"
echo "📊 Health check: http://localhost:$PORT/health"
echo "🔌 API docs: http://localhost:$PORT/api/events"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

./pco-billboard 