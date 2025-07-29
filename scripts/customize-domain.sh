#!/bin/bash

# PCO Arrivals Dashboard - Domain Customization Script
# This script helps customize the configuration for your specific domain

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

log "${BLUE}=== PCO Arrivals Dashboard Domain Customization ===${NC}"

# Check if we're in the right directory
if [ ! -f "env.production" ]; then
    error_exit "Please run this script from the project root directory"
fi

# Get domain information
echo
log "${YELLOW}Please provide your domain information:${NC}"
echo

read -p "Enter your domain name (e.g., mychurch.org): " DOMAIN_NAME
if [ -z "$DOMAIN_NAME" ]; then
    error_exit "Domain name is required"
fi

read -p "Enter your PCO Client ID: " PCO_CLIENT_ID
if [ -z "$PCO_CLIENT_ID" ]; then
    error_exit "PCO Client ID is required"
fi

read -p "Enter your PCO Client Secret: " PCO_CLIENT_SECRET
if [ -z "$PCO_CLIENT_SECRET" ]; then
    error_exit "PCO Client Secret is required"
fi

read -p "Enter your PCO Access Token: " PCO_ACCESS_TOKEN
if [ -z "$PCO_ACCESS_TOKEN" ]; then
    error_exit "PCO Access Token is required"
fi

read -p "Enter your PCO Access Secret: " PCO_ACCESS_SECRET
if [ -z "$PCO_ACCESS_SECRET" ]; then
    error_exit "PCO Access Secret is required"
fi

read -p "Enter your PCO User ID (for admin access): " PCO_USER_ID
if [ -z "$PCO_USER_ID" ]; then
    error_exit "PCO User ID is required"
fi

# Generate secure secrets
SESSION_SECRET=$(openssl rand -base64 64)
JWT_SECRET=$(openssl rand -base64 64)
REDIS_PASSWORD=$(openssl rand -base64 32)

log "${GREEN}Generated secure secrets${NC}"

# Create customized environment file
log "${YELLOW}Creating customized environment file...${NC}"

cat > .env.production << EOF
# Production Environment Configuration
# Customized for domain: $DOMAIN_NAME

# Server Configuration
PORT=3000
HOST=0.0.0.0
ENVIRONMENT=production
TRUST_PROXY=true

# CORS Configuration
CORS_ORIGINS=https://$DOMAIN_NAME,https://www.$DOMAIN_NAME

# Database Configuration
DATABASE_URL=file:./data/pco_billboard.db?cache=shared&mode=rwc&_journal_mode=WAL&_synchronous=NORMAL&_cache_size=10000&_temp_store=MEMORY
DB_MAX_OPEN_CONNS=50
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=600

# PCO OAuth Configuration
PCO_CLIENT_ID=$PCO_CLIENT_ID
PCO_CLIENT_SECRET=$PCO_CLIENT_SECRET
PCO_ACCESS_TOKEN=$PCO_ACCESS_TOKEN
PCO_ACCESS_SECRET=$PCO_ACCESS_SECRET

PCO_REDIRECT_URI=https://$DOMAIN_NAME/auth/callback
PCO_BASE_URL=https://api.planningcenteronline.com
PCO_SCOPES=people check_ins

# Authentication Configuration
SESSION_TTL=3600
REMEMBER_ME_DAYS=30
AUTHORIZED_USERS=$PCO_USER_ID
SESSION_SECRET=$SESSION_SECRET
JWT_SECRET=$JWT_SECRET
TOKEN_REFRESH_THRESHOLD=300

# Redis Configuration
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=$REDIS_PASSWORD
REDIS_DB=0

# Real-time Configuration
WEBSOCKET_ENABLED=true
POLLING_FALLBACK=true
POLLING_INTERVAL=10
LOCATION_POLL_INTERVAL=60
MAX_CONNECTIONS=2000
HEARTBEAT_INTERVAL=30

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
LOG_FILE=/app/logs/app.log
LOG_MAX_SIZE=100
LOG_MAX_BACKUPS=10
LOG_MAX_AGE=30

# Security Configuration
RATE_LIMIT_REQUESTS=200
RATE_LIMIT_WINDOW=60
MAX_REQUEST_SIZE=2097152
SECURITY_HEADERS_ENABLED=true
CSP_ENABLED=true
HSTS_ENABLED=true
HSTS_MAX_AGE=31536000

# Performance Configuration
COMPRESSION_ENABLED=true
CACHE_CONTROL_MAX_AGE=86400
ETAG_ENABLED=true
GZIP_LEVEL=6

# Monitoring Configuration
METRICS_ENABLED=true
HEALTH_CHECK_ENABLED=true
READINESS_CHECK_ENABLED=true
LIVENESS_CHECK_ENABLED=true

# Backup Configuration
BACKUP_ENABLED=true
BACKUP_INTERVAL=24h
BACKUP_RETENTION_DAYS=30
BACKUP_PATH=/app/backups

# SSL/TLS Configuration
SSL_ENABLED=true
SSL_CERT_PATH=/app/ssl/cert.pem
SSL_KEY_PATH=/app/ssl/key.pem
SSL_REDIRECT=true

# Cloudflare Workers Configuration (if using)
CF_WORKERS_ENABLED=false
CF_ACCOUNT_ID=your_cloudflare_account_id
CF_API_TOKEN=your_cloudflare_api_token
CF_ZONE_ID=your_cloudflare_zone_id
EOF

# Update Docker Compose environment variables
log "${YELLOW}Updating Docker Compose configuration...${NC}"

# Create a sed command to update the docker-compose.production.yml
sed -i.bak "s|VITE_API_BASE_URL=\${VITE_API_BASE_URL:-https://your-domain.com}|VITE_API_BASE_URL=\${VITE_API_BASE_URL:-https://$DOMAIN_NAME}|g" docker-compose.production.yml
sed -i.bak "s|VITE_WS_BASE_URL=\${VITE_WS_BASE_URL:-wss://your-domain.com}|VITE_WS_BASE_URL=\${VITE_WS_BASE_URL:-wss://$DOMAIN_NAME}|g" docker-compose.production.yml

# Update Nginx configuration
log "${YELLOW}Updating Nginx configuration...${NC}"

# Create domain-specific Nginx configuration
cat > nginx/conf.d/$DOMAIN_NAME.conf << EOF
server {
    listen 80;
    server_name $DOMAIN_NAME www.$DOMAIN_NAME;
    return 301 https://\$server_name\$request_uri;
}

server {
    listen 443 ssl http2;
    server_name $DOMAIN_NAME www.$DOMAIN_NAME;

    # SSL configuration
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # HSTS
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # API routes
    location /api/ {
        limit_req zone=api burst=20 nodelay;
        
        proxy_pass http://backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
        proxy_read_timeout 86400;
    }

    # WebSocket routes
    location /ws/ {
        proxy_pass http://backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_read_timeout 86400;
    }

    # Health check
    location /health {
        proxy_pass http://backend;
        access_log off;
    }

    # Static files (served by frontend)
    location / {
        proxy_pass http://frontend;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        
        # Cache static assets
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)\$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
            proxy_pass http://frontend;
        }
    }

    # Security: Deny access to sensitive files
    location ~ /\. {
        deny all;
    }

    location ~ /\.(ht|git) {
        deny all;
    }
}
EOF

# Update Cloudflare Workers configuration
log "${YELLOW}Updating Cloudflare Workers configuration...${NC}"

sed -i.bak "s|your-domain.com|$DOMAIN_NAME|g" cloudflare/wrangler.toml
sed -i.bak "s|staging.your-domain.com|staging.$DOMAIN_NAME|g" cloudflare/wrangler.toml

# Create SSL certificate placeholder
log "${YELLOW}Creating SSL certificate placeholder...${NC}"

mkdir -p nginx/ssl
cat > nginx/ssl/README.md << EOF
# SSL Certificates for $DOMAIN_NAME

Place your SSL certificates in this directory:

1. **cert.pem** - Your SSL certificate file
2. **key.pem** - Your SSL private key file

## Getting SSL Certificates

### Option 1: Let's Encrypt (Free)
\`\`\`bash
# Install certbot
sudo apt-get install certbot

# Get certificate
sudo certbot certonly --standalone -d $DOMAIN_NAME -d www.$DOMAIN_NAME

# Copy certificates
sudo cp /etc/letsencrypt/live/$DOMAIN_NAME/fullchain.pem nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/$DOMAIN_NAME/privkey.pem nginx/ssl/key.pem
sudo chown \$USER:\$USER nginx/ssl/*
\`\`\`

### Option 2: Cloudflare (Recommended)
1. Add your domain to Cloudflare
2. Enable SSL/TLS encryption mode to "Full (strict)"
3. Download certificates from Cloudflare dashboard

### Option 3: Commercial Certificate
Purchase from providers like:
- DigiCert
- GlobalSign
- Comodo

## Security
- Set proper permissions: \`chmod 600 nginx/ssl/*\`
- Keep private keys secure
- Renew certificates before expiration
EOF

# Create deployment checklist
log "${YELLOW}Creating deployment checklist...${NC}"

cat > DEPLOYMENT_CHECKLIST.md << EOF
# Deployment Checklist for $DOMAIN_NAME

## ✅ Pre-Deployment Tasks

### 1. SSL Certificates
- [ ] Place SSL certificates in \`nginx/ssl/\`
  - [ ] cert.pem (certificate)
  - [ ] key.pem (private key)
- [ ] Set proper permissions: \`chmod 600 nginx/ssl/*\`

### 2. DNS Configuration
- [ ] Point $DOMAIN_NAME to your server IP
- [ ] Point www.$DOMAIN_NAME to your server IP
- [ ] Verify DNS propagation

### 3. PCO Configuration
- [ ] Verify PCO Client ID: $PCO_CLIENT_ID
- [ ] Verify PCO Client Secret is set
- [ ] Verify PCO Access Token is set
- [ ] Verify PCO Access Secret is set
- [ ] Verify PCO User ID: $PCO_USER_ID

### 4. Server Requirements
- [ ] Docker and Docker Compose installed
- [ ] At least 2GB RAM available
- [ ] At least 20GB disk space
- [ ] Ports 80, 443, 3000 available

## 🚀 Deployment Steps

### 1. Deploy to Production
\`\`\`bash
./scripts/deploy-production.sh
\`\`\`

### 2. Verify Deployment
- [ ] Application accessible at https://$DOMAIN_NAME
- [ ] SSL certificate working
- [ ] PCO authentication working
- [ ] WebSocket connections working
- [ ] Admin panel accessible

### 3. Monitoring Setup
- [ ] Grafana dashboard: http://localhost:3001
- [ ] Prometheus metrics: http://localhost:9090
- [ ] Application health: https://$DOMAIN_NAME/health

## 🔧 Post-Deployment Tasks

### 1. Security
- [ ] Test SSL configuration
- [ ] Verify security headers
- [ ] Test rate limiting
- [ ] Verify CORS settings

### 2. Performance
- [ ] Monitor response times
- [ ] Check memory usage
- [ ] Verify caching
- [ ] Test WebSocket connections

### 3. Backup
- [ ] Test backup script: \`./scripts/backup.sh\`
- [ ] Set up automated backups
- [ ] Verify backup restoration

## 📞 Support Information

- **Application URL**: https://$DOMAIN_NAME
- **Admin Panel**: https://$DOMAIN_NAME/admin
- **Billboard**: https://$DOMAIN_NAME/billboard
- **Monitoring**: http://localhost:3001 (Grafana)

## 🔒 Security Notes

- Keep \`.env.production\` secure
- Regularly update SSL certificates
- Monitor logs for security issues
- Set up automated security updates
EOF

# Create a quick test script
log "${YELLOW}Creating quick test script...${NC}"

cat > scripts/test-domain.sh << 'EOF'
#!/bin/bash

# Quick domain test script
DOMAIN_NAME=$(grep "CORS_ORIGINS" .env.production | cut -d'=' -f2 | cut -d',' -f1 | sed 's|https://||')

echo "Testing domain configuration for: $DOMAIN_NAME"

# Test DNS resolution
echo "Testing DNS resolution..."
if nslookup $DOMAIN_NAME > /dev/null 2>&1; then
    echo "✅ DNS resolution working"
else
    echo "❌ DNS resolution failed"
fi

# Test SSL certificate (if available)
if [ -f "nginx/ssl/cert.pem" ]; then
    echo "✅ SSL certificate found"
else
    echo "⚠️  SSL certificate not found - please add certificates to nginx/ssl/"
fi

# Test Docker configuration
echo "Testing Docker configuration..."
if docker-compose -f docker-compose.production.yml config > /dev/null 2>&1; then
    echo "✅ Docker Compose configuration valid"
else
    echo "❌ Docker Compose configuration invalid"
fi

echo "Domain configuration test complete!"
EOF

chmod +x scripts/test-domain.sh

# Summary
log "${GREEN}=== Domain Customization Complete! ===${NC}"
echo
log "${BLUE}Configuration Summary:${NC}"
echo "  Domain: $DOMAIN_NAME"
echo "  PCO Client ID: $PCO_CLIENT_ID"
echo "  PCO User ID: $PCO_USER_ID"
echo "  Application URL: https://$DOMAIN_NAME"
echo "  Admin Panel: https://$DOMAIN_NAME/admin"
echo "  Billboard: https://$DOMAIN_NAME/billboard"
echo
log "${YELLOW}Next Steps:${NC}"
echo "  1. Add SSL certificates to nginx/ssl/"
echo "  2. Configure DNS for $DOMAIN_NAME"
echo "  3. Run: ./scripts/test-domain.sh"
echo "  4. Run: ./scripts/deploy-production.sh"
echo
log "${GREEN}Files Created/Updated:${NC}"
echo "  ✅ .env.production (customized)"
echo "  ✅ docker-compose.production.yml (updated)"
echo "  ✅ nginx/conf.d/$DOMAIN_NAME.conf"
echo "  ✅ cloudflare/wrangler.toml (updated)"
echo "  ✅ nginx/ssl/README.md"
echo "  ✅ DEPLOYMENT_CHECKLIST.md"
echo "  ✅ scripts/test-domain.sh"
echo
log "${GREEN}Domain customization completed successfully! 🚀${NC}" 