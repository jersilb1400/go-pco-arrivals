#!/bin/bash

# Railway Setup Script for PCO Arrivals Dashboard
# This script automates the setup process for Railway deployment

set -e

echo "🚂 Railway Deployment Setup"
echo "============================"

# Check if Railway CLI is installed
if ! command -v railway &> /dev/null; then
    echo "❌ Railway CLI not found. Installing..."
    npm install -g @railway/cli
else
    echo "✅ Railway CLI found"
fi

# Check if user is logged in
if ! railway whoami &> /dev/null; then
    echo "🔐 Please login to Railway..."
    railway login
else
    echo "✅ Already logged in to Railway"
fi

echo ""
echo "📋 Setup Steps:"
echo "1. Initialize Railway project"
echo "2. Configure environment variables"
echo "3. Deploy the application"
echo "4. Set up custom domain (optional)"
echo ""

# Function to initialize Railway project
init_project() {
    echo "🚀 Initializing Railway project..."
    
    # Initialize Railway project
    railway init
    
    echo "✅ Railway project initialized"
}

# Function to set environment variables
set_env_vars() {
    echo "🔧 Setting up environment variables..."
    
    echo "Setting production environment variables..."
    
    # Required environment variables
    railway variables set ENVIRONMENT=production
    railway variables set PORT=3000
    railway variables set HOST=0.0.0.0
    railway variables set TRUST_PROXY=true
    railway variables set DATABASE_TYPE=mongodb
    railway variables set MONGODB_DATABASE=go-pco-arrivals
    railway variables set LOG_LEVEL=info
    railway variables set METRICS_ENABLED=true
    
    echo "✅ Environment variables set"
    echo ""
    echo "⚠️  IMPORTANT: You need to set these secrets manually:"
    echo "   - MONGODB_URI (MongoDB Atlas connection string)"
    echo "   - PCO_CLIENT_ID (Planning Center client ID)"
    echo "   - PCO_CLIENT_SECRET (Planning Center client secret)"
    echo "   - SESSION_SECRET (Session encryption key)"
    echo "   - JWT_SECRET (JWT signing key)"
    echo ""
    echo "To set secrets, run:"
    echo "  railway variables set MONGODB_URI='your-mongodb-uri'"
    echo "  railway variables set PCO_CLIENT_ID='your-pco-client-id'"
    echo "  railway variables set PCO_CLIENT_SECRET='your-pco-client-secret'"
    echo "  railway variables set SESSION_SECRET='your-session-secret'"
    echo "  railway variables set JWT_SECRET='your-jwt-secret'"
}

# Function to deploy the application
deploy_app() {
    echo "🚀 Deploying to Railway..."
    
    # Build the application
    echo "Building Go application..."
    go build -o main .
    
    if [ $? -eq 0 ]; then
        echo "✅ Build successful"
        
        # Deploy to Railway
        echo "Deploying to Railway..."
        railway up
        
        if [ $? -eq 0 ]; then
            echo "✅ Deployment successful!"
            echo "🌐 Your app is now live!"
            
            # Get the deployment URL
            echo "🔗 Getting deployment URL..."
            railway status
        else
            echo "❌ Deployment failed"
            exit 1
        fi
    else
        echo "❌ Build failed"
        exit 1
    fi
}

# Function to show deployment status
show_status() {
    echo ""
    echo "📊 Deployment Status:"
    echo "===================="
    
    railway status
    
    echo ""
    echo "🔍 Check your deployment:"
    echo "- Dashboard: https://railway.app/dashboard"
    echo "- Logs: railway logs"
    echo "- Status: railway status"
}

# Function to set secrets interactively
set_secrets() {
    echo "🔐 Setting up secrets..."
    
    echo "Enter your MongoDB Atlas connection string:"
    read -s MONGODB_URI
    railway variables set MONGODB_URI="$MONGODB_URI"
    
    echo "Enter your PCO Client ID:"
    read -s PCO_CLIENT_ID
    railway variables set PCO_CLIENT_ID="$PCO_CLIENT_ID"
    
    echo "Enter your PCO Client Secret:"
    read -s PCO_CLIENT_SECRET
    railway variables set PCO_CLIENT_SECRET="$PCO_CLIENT_SECRET"
    
    echo "Enter your Session Secret (or press Enter to generate):"
    read -s SESSION_SECRET
    if [ -z "$SESSION_SECRET" ]; then
        SESSION_SECRET=$(openssl rand -base64 32)
        echo "Generated Session Secret: $SESSION_SECRET"
    fi
    railway variables set SESSION_SECRET="$SESSION_SECRET"
    
    echo "Enter your JWT Secret (or press Enter to generate):"
    read -s JWT_SECRET
    if [ -z "$JWT_SECRET" ]; then
        JWT_SECRET=$(openssl rand -base64 32)
        echo "Generated JWT Secret: $JWT_SECRET"
    fi
    railway variables set JWT_SECRET="$JWT_SECRET"
    
    echo "✅ Secrets configured successfully!"
}

# Main menu
while true; do
    echo ""
    echo "Choose an option:"
    echo "1) Initialize Railway project"
    echo "2) Set environment variables"
    echo "3) Set secrets interactively"
    echo "4) Deploy application"
    echo "5) Deploy everything"
    echo "6) Show deployment status"
    echo "7) Exit"
    echo ""
    read -p "Enter your choice (1-7): " choice
    
    case $choice in
        1)
            init_project
            ;;
        2)
            set_env_vars
            ;;
        3)
            set_secrets
            ;;
        4)
            deploy_app
            ;;
        5)
            init_project
            set_env_vars
            set_secrets
            deploy_app
            show_status
            ;;
        6)
            show_status
            ;;
        7)
            echo "👋 Goodbye!"
            exit 0
            ;;
        *)
            echo "❌ Invalid choice. Please try again."
            ;;
    esac
done 