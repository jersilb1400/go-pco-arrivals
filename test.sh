#!/bin/bash

# PCO Arrivals Billboard - Test Script

echo "🧪 Testing PCO Arrivals Billboard..."

# Wait for server to start
sleep 2

# Test health endpoint
echo "📊 Testing health endpoint..."
curl -s http://localhost:3000/health | jq '.' 2>/dev/null || curl -s http://localhost:3000/health

echo ""
echo "📊 Testing detailed health endpoint..."
curl -s http://localhost:3000/health/detailed | jq '.' 2>/dev/null || curl -s http://localhost:3000/health/detailed

echo ""
echo "🔐 Testing auth status endpoint..."
curl -s http://localhost:3000/auth/status | jq '.' 2>/dev/null || curl -s http://localhost:3000/auth/status

echo ""
echo "📅 Testing events endpoint..."
curl -s http://localhost:3000/api/events | jq '.' 2>/dev/null || curl -s http://localhost:3000/api/events

echo ""
echo "📺 Testing billboard state endpoint..."
curl -s http://localhost:3000/billboard/state | jq '.' 2>/dev/null || curl -s http://localhost:3000/billboard/state

echo ""
echo "✅ All tests completed!" 