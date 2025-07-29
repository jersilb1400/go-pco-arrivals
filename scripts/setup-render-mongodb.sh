#!/bin/bash

# PCO Arrivals Dashboard - Render + MongoDB Atlas Setup Script
# This script helps configure the application for Render hosting with MongoDB Atlas

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "[$(date +'%Y-%m-%d %H:%M:%S')] $1"
}

# Error handling
error_exit() {
    log "${RED}ERROR: $1${NC}" >&2
    exit 1
}

log "${BLUE}=== PCO Arrivals Dashboard - Render + MongoDB Atlas Setup ===${NC}"

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    error_exit "Please run this script from the project root directory"
fi

# Get Render and MongoDB Atlas information
echo
log "${YELLOW}Please provide your Render and MongoDB Atlas information:${NC}"
echo

read -p "Enter your Render service name (e.g., pco-arrivals-dashboard): " RENDER_SERVICE_NAME
read -p "Enter your Render service URL (e.g., https://your-app.onrender.com): " RENDER_SERVICE_URL
read -p "Enter your MongoDB Atlas connection string: " MONGODB_ATLAS_URI
read -p "Enter your MongoDB Atlas database name: " MONGODB_DATABASE_NAME

# Get PCO credentials
echo
log "${YELLOW}PCO (Planning Center Online) Configuration:${NC}"
echo

read -p "Enter your PCO Client ID: " PCO_CLIENT_ID
read -p "Enter your PCO Client Secret: " PCO_CLIENT_SECRET

# Generate secure secrets
SESSION_SECRET=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 32)

# Create Render-specific environment file
log "${GREEN}Creating Render environment configuration...${NC}"

cat > .env.render << EOF
# Render Production Environment Configuration
# ========================================

# Application Settings
ENVIRONMENT=production
PORT=3000
HOST=0.0.0.0
TRUST_PROXY=true

# CORS Configuration
CORS_ORIGINS=${RENDER_SERVICE_URL}
CORS_ALLOW_CREDENTIALS=true
CORS_ALLOW_HEADERS=Content-Type,Authorization,X-Requested-With
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS

# MongoDB Atlas Configuration
DATABASE_TYPE=mongodb
MONGODB_URI=${MONGODB_ATLAS_URI}
MONGODB_DATABASE=${MONGODB_DATABASE_NAME}
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=300s

# PCO Configuration
PCO_CLIENT_ID=${PCO_CLIENT_ID}
PCO_CLIENT_SECRET=${PCO_CLIENT_SECRET}
PCO_REDIRECT_URI=${RENDER_SERVICE_URL}/auth/callback

# Security Configuration
SESSION_SECRET=${SESSION_SECRET}
JWT_SECRET=${JWT_SECRET}
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
MAX_REQUEST_SIZE=10MB
SECURITY_HEADERS_ENABLED=true
CSP_ENABLED=true
HSTS_ENABLED=true
HSTS_MAX_AGE=31536000

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
LOG_FILE=/tmp/app.log

# Performance Configuration
GZIP_LEVEL=6
COMPRESSION_ENABLED=true

# Monitoring Configuration
METRICS_ENABLED=true
HEALTH_CHECK_ENABLED=true
READINESS_CHECK_ENABLED=true
LIVENESS_CHECK_ENABLED=true

# WebSocket Configuration
MAX_CONNECTIONS=1000
HEARTBEAT_INTERVAL=30s
EOF

# Create render.yaml for Render deployment
log "${GREEN}Creating Render deployment configuration...${NC}"

cat > render.yaml << EOF
services:
  - type: web
    name: ${RENDER_SERVICE_NAME}
    env: docker
    plan: starter
    region: oregon
    buildCommand: |
      # Build the Go backend
      go mod download
      go build -o main .
      
      # Build the React frontend
      cd frontend
      npm ci
      npm run build
      cd ..
      
      # Copy frontend build to backend
      cp -r frontend/dist web/static
    startCommand: ./main
    envVars:
      - key: ENVIRONMENT
        value: production
      - key: PORT
        value: 3000
      - key: HOST
        value: 0.0.0.0
      - key: TRUST_PROXY
        value: true
      - key: CORS_ORIGINS
        value: ${RENDER_SERVICE_URL}
      - key: DATABASE_TYPE
        value: mongodb
      - key: MONGODB_URI
        sync: false
      - key: MONGODB_DATABASE
        value: ${MONGODB_DATABASE_NAME}
      - key: PCO_CLIENT_ID
        sync: false
      - key: PCO_CLIENT_SECRET
        sync: false
      - key: PCO_REDIRECT_URI
        value: ${RENDER_SERVICE_URL}/auth/callback
      - key: SESSION_SECRET
        sync: false
      - key: JWT_SECRET
        sync: false
      - key: LOG_LEVEL
        value: info
      - key: METRICS_ENABLED
        value: true
    healthCheckPath: /health
    autoDeploy: true
    branch: main

databases:
  - name: ${MONGODB_DATABASE_NAME}
    databaseName: ${MONGODB_DATABASE_NAME}
    plan: free
    region: oregon
EOF

# Create Dockerfile for Render
log "${GREEN}Creating Render-optimized Dockerfile...${NC}"

cat > Dockerfile.render << EOF
# Multi-stage build for Render deployment
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata build-base nodejs npm

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the Go backend
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags="-w -s -X main.Version=\$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
    -o main .

# Build the React frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci --only=production
COPY frontend/ .
RUN npm run build

# Production stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Create necessary directories
RUN mkdir -p /app/data /app/logs /app/web/static && \
    chown -R appuser:appgroup /app

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Copy frontend build
COPY --from=builder /app/frontend/dist ./web/static

# Set proper ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \\
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1

# Run the application
CMD ["./main"]
EOF

# Create MongoDB Atlas setup guide
log "${GREEN}Creating MongoDB Atlas setup guide...${NC}"

cat > MONGODB_ATLAS_SETUP.md << EOF
# MongoDB Atlas Setup Guide

## 1. Create MongoDB Atlas Cluster

1. Go to [MongoDB Atlas](https://cloud.mongodb.com)
2. Sign in to your account
3. Click "Build a Database"
4. Choose "FREE" tier (M0)
5. Select your preferred cloud provider and region
6. Click "Create"

## 2. Configure Database Access

1. Go to "Database Access" in the left sidebar
2. Click "Add New Database User"
3. Choose "Password" authentication
4. Create a username and password
5. Set privileges to "Read and write to any database"
6. Click "Add User"

## 3. Configure Network Access

1. Go to "Network Access" in the left sidebar
2. Click "Add IP Address"
3. For Render deployment, add: \`0.0.0.0/0\` (allows all IPs)
4. Click "Confirm"

## 4. Get Connection String

1. Go to "Database" in the left sidebar
2. Click "Connect"
3. Choose "Connect your application"
4. Copy the connection string
5. Replace \`<password>\` with your database user password
6. Replace \`<dbname>\` with your database name

## 5. Update Render Environment Variables

In your Render dashboard:
1. Go to your service
2. Click "Environment"
3. Add the following variables:
   - \`MONGODB_URI\`: Your connection string
   - \`MONGODB_DATABASE\`: Your database name
   - \`PCO_CLIENT_ID\`: Your PCO Client ID
   - \`PCO_CLIENT_SECRET\`: Your PCO Client Secret
   - \`SESSION_SECRET\`: Generated secret
   - \`JWT_SECRET\`: Generated secret

## 6. Deploy to Render

1. Push your code to GitHub
2. Connect your GitHub repository to Render
3. Render will automatically deploy using the \`render.yaml\` configuration
4. Your app will be available at your Render URL

## 7. Test the Deployment

1. Visit your Render URL
2. Test the login functionality
3. Verify WebSocket connections work
4. Check that data is being stored in MongoDB Atlas

## Troubleshooting

- **Connection Issues**: Ensure your MongoDB Atlas cluster is in the same region as your Render service
- **Authentication Errors**: Double-check your database username and password
- **Network Access**: Make sure you've added \`0.0.0.0/0\` to allowed IP addresses
- **Environment Variables**: Verify all required environment variables are set in Render
EOF

# Create deployment checklist
log "${GREEN}Creating deployment checklist...${NC}"

cat > RENDER_DEPLOYMENT_CHECKLIST.md << EOF
# Render Deployment Checklist

## ✅ Pre-Deployment Checklist

### MongoDB Atlas Setup
- [ ] MongoDB Atlas cluster created
- [ ] Database user created with read/write permissions
- [ ] Network access configured (0.0.0.0/0)
- [ ] Connection string obtained and tested

### PCO Configuration
- [ ] PCO application created in Planning Center
- [ ] Client ID and Secret obtained
- [ ] Redirect URI configured: \`${RENDER_SERVICE_URL}/auth/callback\`
- [ ] Required scopes configured

### Render Configuration
- [ ] Render account created
- [ ] GitHub repository connected to Render
- [ ] Environment variables configured in Render dashboard
- [ ] Auto-deploy enabled

### Code Preparation
- [ ] All changes committed to main branch
- [ ] render.yaml file present in repository
- [ ] Dockerfile.render present in repository
- [ ] .env.render file created (for reference)

## 🚀 Deployment Steps

1. **Push to GitHub**
   \`\`\`bash
   git add .
   git commit -m "Configure for Render deployment"
   git push origin main
   \`\`\`

2. **Deploy on Render**
   - Render will automatically detect the render.yaml file
   - Build will start automatically
   - Monitor the build logs for any issues

3. **Configure Environment Variables**
   In Render dashboard, set these variables:
   - \`MONGODB_URI\`: Your MongoDB Atlas connection string
   - \`MONGODB_DATABASE\`: Your database name
   - \`PCO_CLIENT_ID\`: Your PCO Client ID
   - \`PCO_CLIENT_SECRET\`: Your PCO Client Secret
   - \`SESSION_SECRET\`: ${SESSION_SECRET}
   - \`JWT_SECRET\`: ${JWT_SECRET}

4. **Test the Deployment**
   - Visit your Render URL
   - Test login functionality
   - Verify WebSocket connections
   - Check MongoDB Atlas for data storage

## 🔧 Post-Deployment Verification

- [ ] Application loads without errors
- [ ] PCO login works correctly
- [ ] WebSocket connections establish
- [ ] Real-time updates work
- [ ] Data is being stored in MongoDB Atlas
- [ ] Health check endpoint responds
- [ ] SSL/HTTPS is working

## 📊 Monitoring

- Monitor Render logs for any errors
- Check MongoDB Atlas metrics
- Verify application performance
- Test all features thoroughly

## 🆘 Troubleshooting

### Common Issues:
1. **Build Failures**: Check render.yaml syntax
2. **Connection Errors**: Verify MongoDB Atlas settings
3. **Authentication Issues**: Check PCO credentials
4. **Environment Variables**: Ensure all required vars are set

### Support Resources:
- Render Documentation: https://render.com/docs
- MongoDB Atlas Documentation: https://docs.atlas.mongodb.com
- PCO API Documentation: https://developer.planningcenteronline.com
EOF

# Make scripts executable
chmod +x scripts/*.sh

log "${GREEN}✅ Render + MongoDB Atlas setup complete!${NC}"
echo
log "${YELLOW}📋 Next Steps:${NC}"
echo
log "1. Review the generated files:"
log "   - .env.render (environment configuration)"
log "   - render.yaml (Render deployment config)"
log "   - Dockerfile.render (Render-optimized Dockerfile)"
log "   - MONGODB_ATLAS_SETUP.md (MongoDB setup guide)"
log "   - RENDER_DEPLOYMENT_CHECKLIST.md (deployment checklist)"
echo
log "2. Set up MongoDB Atlas:"
log "   - Follow the instructions in MONGODB_ATLAS_SETUP.md"
log "   - Get your connection string"
echo
log "3. Configure Render:"
log "   - Push your code to GitHub"
log "   - Connect your repository to Render"
log "   - Set environment variables in Render dashboard"
echo
log "4. Deploy:"
log "   - Render will automatically deploy using render.yaml"
log "   - Monitor the build logs"
log "   - Test all functionality"
echo
log "${GREEN}🎉 Your PCO Arrivals Dashboard is ready for Render deployment!${NC}" 