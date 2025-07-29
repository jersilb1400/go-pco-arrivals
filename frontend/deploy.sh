#!/bin/bash

# PCO Arrivals Billboard Frontend Deployment Script

set -e

echo "🚀 Starting frontend deployment..."

# Check if we're in the right directory
if [ ! -f "package.json" ]; then
    echo "❌ Error: package.json not found. Please run this script from the frontend directory."
    exit 1
fi

# Install dependencies
echo "📦 Installing dependencies..."
npm ci --only=production

# Build the application
echo "🔨 Building application..."
npm run build

# Check if build was successful
if [ ! -d "dist" ]; then
    echo "❌ Error: Build failed - dist directory not found"
    exit 1
fi

echo "✅ Build completed successfully!"

# Optional: Deploy to a static hosting service
# Uncomment and configure as needed

# # Deploy to Netlify (if netlify-cli is installed)
# if command -v netlify &> /dev/null; then
#     echo "🌐 Deploying to Netlify..."
#     netlify deploy --prod --dir=dist
# fi

# # Deploy to Vercel (if vercel is installed)
# if command -v vercel &> /dev/null; then
#     echo "🌐 Deploying to Vercel..."
#     vercel --prod
# fi

echo "🎉 Deployment completed!"
echo "📁 Build output is in the 'dist' directory"
echo "🌐 You can now deploy the contents of 'dist' to your hosting provider" 