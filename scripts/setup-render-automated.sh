#!/bin/bash

# PCO Arrivals Dashboard - Automated Render Setup Script
# This script helps automate the Render deployment setup process

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
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

# Success message
success() {
    log "${GREEN}✅ $1${NC}"
}

# Info message
info() {
    log "${BLUE}ℹ️  $1${NC}"
}

# Warning message
warning() {
    log "${YELLOW}⚠️  $1${NC}"
}

log "${PURPLE}=== PCO Arrivals Dashboard - Automated Render Setup ===${NC}"

# Check if we're in the right directory
if [ ! -f "render.yaml" ]; then
    error_exit "Please run this script from the project root directory"
fi

# Check if GitHub repository exists
if ! git remote get-url origin >/dev/null 2>&1; then
    error_exit "Git repository not configured. Please run 'git remote add origin' first."
fi

REPO_URL=$(git remote get-url origin)
REPO_NAME=$(basename -s .git "$REPO_URL")

info "Repository detected: $REPO_NAME"
info "Repository URL: $REPO_URL"

# Configuration values
RENDER_SERVICE_NAME="go-pco-arrivals-dashboard"
RENDER_SERVICE_URL="https://go-pco-arrivals-dashboard.onrender.com"
MONGODB_URI="mongodb+srv://jeremy:tn8KBdz5M39nhCNr@cluster0.8zrrvsj.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
MONGODB_DATABASE="go-pco-arrivals-dashboard"
PCO_CLIENT_ID="YOUR_PCO_CLIENT_ID_HERE"
PCO_CLIENT_SECRET="YOUR_PCO_CLIENT_SECRET_HERE"
SESSION_SECRET="YIVWPwOmXOr8UY3cwXD7eJeVGFiL/vuV6/gLdrOBnzA="
JWT_SECRET="88Os98iSRtVVtHnAYHvPt6c9eatogKIK2BVGum4l74M="

echo
log "${CYAN}🚀 Starting Automated Render Setup...${NC}"
echo

# Step 1: Verify GitHub repository is ready
info "Step 1: Verifying GitHub repository..."
if git status --porcelain | grep -q .; then
    warning "You have uncommitted changes. Consider committing them first:"
    echo "  git add . && git commit -m 'Update for Render deployment'"
    echo "  git push origin main"
    echo
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    success "Repository is clean and ready for deployment"
fi

# Step 2: Create deployment checklist
info "Step 2: Creating deployment checklist..."
cat > RENDER_AUTOMATED_SETUP.md << EOF
# 🚀 Automated Render Setup Checklist

## ✅ Pre-Deployment Verification

### GitHub Repository
- [x] Repository: \`$REPO_NAME\`
- [x] URL: \`$REPO_URL\`
- [x] Branch: \`main\`
- [x] render.yaml: ✅ Present
- [x] Dockerfile.render: ✅ Present

### Configuration Files
- [x] render.yaml: ✅ Configured
- [x] Dockerfile.render: ✅ Optimized for Render
- [x] .env.render: ✅ Environment template
- [x] MONGODB_ATLAS_SETUP.md: ✅ Setup guide

## 🌐 Render Dashboard Setup

### 1. Create New Web Service
1. Go to [Render Dashboard](https://dashboard.render.com)
2. Click "New +" → "Web Service"
3. Connect GitHub account (if not already connected)
4. Select repository: \`$REPO_NAME\`

### 2. Service Configuration
Render will automatically detect and configure:
- **Name**: \`$RENDER_SERVICE_NAME\`
- **Environment**: Docker
- **Build Command**: Automatically set from render.yaml
- **Start Command**: \`./main\`
- **Health Check Path**: \`/health\`

### 3. Environment Variables
Add these environment variables in Render dashboard:

#### Required Variables:
\`\`\`
MONGODB_URI=$MONGODB_URI
MONGODB_DATABASE=$MONGODB_DATABASE
PCO_CLIENT_ID=$PCO_CLIENT_ID
PCO_CLIENT_SECRET=$PCO_CLIENT_SECRET
SESSION_SECRET=$SESSION_SECRET
JWT_SECRET=$JWT_SECRET
\`\`\`

#### Optional Variables (already in render.yaml):
\`\`\`
ENVIRONMENT=production
PORT=3000
HOST=0.0.0.0
TRUST_PROXY=true
CORS_ORIGINS=$RENDER_SERVICE_URL
DATABASE_TYPE=mongodb
PCO_REDIRECT_URI=$RENDER_SERVICE_URL/auth/callback
LOG_LEVEL=info
METRICS_ENABLED=true
\`\`\`

### 4. PCO Configuration Update
Update your Planning Center application settings:
- **Redirect URI**: \`$RENDER_SERVICE_URL/auth/callback\`

## 🔧 Deployment Steps

### 1. Deploy on Render
1. Click "Create Web Service"
2. Monitor build logs for any issues
3. Wait for deployment to complete

### 2. Test the Deployment
1. Visit: \`$RENDER_SERVICE_URL\`
2. Test PCO login functionality
3. Verify WebSocket connections
4. Check MongoDB Atlas for data storage

### 3. Monitor and Debug
- Check Render logs for any errors
- Verify all environment variables are set
- Test all application features

## 📊 Post-Deployment Verification

- [ ] Application loads without errors
- [ ] PCO login works correctly
- [ ] WebSocket connections establish
- [ ] Real-time updates work
- [ ] Data is being stored in MongoDB Atlas
- [ ] Health check endpoint responds
- [ ] SSL/HTTPS is working

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

## 🔗 Quick Links

- **Render Dashboard**: https://dashboard.render.com
- **GitHub Repository**: $REPO_URL
- **Application URL**: $RENDER_SERVICE_URL
- **MongoDB Atlas**: https://cloud.mongodb.com
- **Planning Center**: https://planningcenteronline.com

## 📝 Environment Variables Summary

\`\`\`bash
# Copy these to Render Environment Variables
MONGODB_URI=$MONGODB_URI
MONGODB_DATABASE=$MONGODB_DATABASE
PCO_CLIENT_ID=$PCO_CLIENT_ID
PCO_CLIENT_SECRET=$PCO_CLIENT_SECRET
SESSION_SECRET=$SESSION_SECRET
JWT_SECRET=$JWT_SECRET
\`\`\`
EOF

success "Created RENDER_AUTOMATED_SETUP.md with complete deployment guide"

# Step 3: Create environment variables file for easy copying
info "Step 3: Creating environment variables file..."
cat > RENDER_ENV_VARS.txt << EOF
# Copy these environment variables to Render Dashboard

MONGODB_URI=$MONGODB_URI
MONGODB_DATABASE=$MONGODB_DATABASE
PCO_CLIENT_ID=$PCO_CLIENT_ID
PCO_CLIENT_SECRET=$PCO_CLIENT_SECRET
SESSION_SECRET=$SESSION_SECRET
JWT_SECRET=$JWT_SECRET
EOF

success "Created RENDER_ENV_VARS.txt with environment variables"

# Step 4: Create quick setup script
info "Step 4: Creating quick setup script..."
cat > scripts/quick-render-setup.sh << 'EOF'
#!/bin/bash

# Quick Render Setup Helper
echo "🚀 Quick Render Setup Helper"
echo "=============================="
echo
echo "1. Go to: https://dashboard.render.com"
echo "2. Click 'New +' → 'Web Service'"
echo "3. Select repository: $(basename -s .git "$(git remote get-url origin)")"
echo "4. Add these environment variables:"
echo
cat RENDER_ENV_VARS.txt
echo
echo "5. Click 'Create Web Service'"
echo "6. Monitor build logs"
echo "7. Test at: https://go-pco-arrivals-dashboard.onrender.com"
echo
echo "📋 Full guide: RENDER_AUTOMATED_SETUP.md"
EOF

chmod +x scripts/quick-render-setup.sh
success "Created quick setup script: scripts/quick-render-setup.sh"

# Step 5: Verify all required files exist
info "Step 5: Verifying required files..."
required_files=(
    "render.yaml"
    "Dockerfile.render"
    ".env.render"
    "MONGODB_ATLAS_SETUP.md"
    "RENDER_DEPLOYMENT_CHECKLIST.md"
)

for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        success "✅ $file"
    else
        error_exit "❌ Missing required file: $file"
    fi
done

# Step 6: Create deployment status checker
info "Step 6: Creating deployment status checker..."
cat > scripts/check-deployment.sh << 'EOF'
#!/bin/bash

# Check deployment status
echo "🔍 Checking deployment status..."
echo "================================"

# Check if the service is responding
URL="https://go-pco-arrivals-dashboard.onrender.com"

echo "Testing application health..."
if curl -s -f "$URL/health" > /dev/null; then
    echo "✅ Application is running"
else
    echo "❌ Application is not responding"
fi

echo
echo "Testing main page..."
if curl -s -f "$URL" > /dev/null; then
    echo "✅ Main page is accessible"
else
    echo "❌ Main page is not accessible"
fi

echo
echo "🔗 Application URL: $URL"
echo "📊 Render Dashboard: https://dashboard.render.com"
EOF

chmod +x scripts/check-deployment.sh
success "Created deployment checker: scripts/check-deployment.sh"

# Step 7: Final summary
echo
log "${GREEN}🎉 Automated Render Setup Complete!${NC}"
echo
log "${CYAN}📋 Next Steps:${NC}"
echo
log "1. 📖 Review the setup guide:"
log "   cat RENDER_AUTOMATED_SETUP.md"
echo
log "2. 🚀 Run the quick setup helper:"
log "   ./scripts/quick-render-setup.sh"
echo
log "3. 🌐 Go to Render Dashboard:"
log "   https://dashboard.render.com"
echo
log "4. 📝 Copy environment variables from:"
log "   cat RENDER_ENV_VARS.txt"
echo
log "5. ✅ Check deployment status:"
log "   ./scripts/check-deployment.sh"
echo
log "${PURPLE}🎯 Your application will be available at:${NC}"
log "${GREEN}$RENDER_SERVICE_URL${NC}"
echo
log "${YELLOW}💡 Pro Tips:${NC}"
log "• Monitor build logs in Render dashboard"
log "• Test all features after deployment"
log "• Check MongoDB Atlas for data storage"
log "• Verify PCO redirect URI is updated"
echo
success "Setup complete! Ready for Render deployment! 🚀" 