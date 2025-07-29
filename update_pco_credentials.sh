#!/bin/bash

echo "🔧 PCO OAuth App Credentials Update Tool"
echo "========================================"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Please provide your new PCO OAuth app credentials:${NC}"
echo ""

# Get new credentials
read -p "New Client ID: " new_client_id
read -p "New Client Secret: " new_client_secret

if [ -z "$new_client_id" ] || [ -z "$new_client_secret" ]; then
    echo -e "${YELLOW}⚠️  Credentials cannot be empty. Please try again.${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}Updating .env file...${NC}"

# Check if .env exists
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}⚠️  .env file not found. Creating from template...${NC}"
    cp env.example .env
fi

# Update the credentials
sed -i.bak "s/PCO_CLIENT_ID=.*/PCO_CLIENT_ID=$new_client_id/" .env
sed -i.bak "s/PCO_CLIENT_SECRET=.*/PCO_CLIENT_SECRET=$new_client_secret/" .env

# Verify the update
echo ""
echo -e "${GREEN}✅ Credentials updated successfully!${NC}"
echo ""
echo "Updated values:"
echo "Client ID: ${new_client_id:0:8}..."
echo "Client Secret: ${new_client_secret:0:8}..."

echo ""
echo -e "${BLUE}Next steps:${NC}"
echo "1. Restart your backend server"
echo "2. Log out and log back in to get a fresh token"
echo "3. Test the APIs with: ./test_pco_status.sh"

# Clean up backup file
rm -f .env.bak 