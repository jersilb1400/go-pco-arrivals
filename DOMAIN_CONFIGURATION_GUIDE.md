# 🌐 Domain Configuration Guide

This guide will help you customize your PCO Arrivals Dashboard for your specific domain.

## 🚀 Quick Start

### Option 1: Automated Configuration (Recommended)

Run the automated domain customization script:

```bash
./scripts/customize-domain.sh
```

This script will:
- ✅ Prompt for your domain and PCO credentials
- ✅ Generate secure secrets automatically
- ✅ Update all configuration files
- ✅ Create SSL certificate instructions
- ✅ Generate a deployment checklist

### Option 2: Manual Configuration

If you prefer to configure manually, follow the steps below.

## 📋 Required Information

Before starting, gather the following information:

### Domain Information
- **Domain Name**: Your domain (e.g., `mychurch.org`)
- **SSL Certificates**: SSL certificate and private key files

### PCO API Credentials
- **PCO Client ID**: From your PCO application
- **PCO Client Secret**: From your PCO application  
- **PCO Access Token**: Your PCO access token
- **PCO Access Secret**: Your PCO access secret
- **PCO User ID**: Your PCO user ID for admin access

## 🔧 Manual Configuration Steps

### Step 1: Update Environment Variables

Edit `.env.production` and update these values:

```bash
# CORS Configuration
CORS_ORIGINS=https://your-domain.com,https://www.your-domain.com

# PCO OAuth Configuration
PCO_CLIENT_ID=your_actual_pco_client_id
PCO_CLIENT_SECRET=your_actual_pco_client_secret
PCO_ACCESS_TOKEN=your_actual_pco_access_token
PCO_ACCESS_SECRET=your_actual_pco_access_secret

PCO_REDIRECT_URI=https://your-domain.com/auth/callback

# Authentication Configuration
AUTHORIZED_USERS=your_actual_pco_user_id

# Generate secure secrets
SESSION_SECRET=$(openssl rand -base64 64)
JWT_SECRET=$(openssl rand -base64 64)
REDIS_PASSWORD=$(openssl rand -base64 32)
```

### Step 2: Update Docker Compose Configuration

Edit `docker-compose.production.yml` and update:

```yaml
environment:
  - VITE_API_BASE_URL=${VITE_API_BASE_URL:-https://your-domain.com}
  - VITE_WS_BASE_URL=${VITE_WS_BASE_URL:-wss://your-domain.com}
```

### Step 3: Configure SSL Certificates

1. **Create SSL directory**:
   ```bash
   mkdir -p nginx/ssl
   ```

2. **Add your SSL certificates**:
   ```bash
   cp your-certificate.pem nginx/ssl/cert.pem
   cp your-private-key.pem nginx/ssl/key.pem
   chmod 600 nginx/ssl/*
   ```

### Step 4: Update Nginx Configuration

Create `nginx/conf.d/your-domain.conf`:

```nginx
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com www.your-domain.com;

    # SSL configuration
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    # ... rest of configuration
}
```

### Step 5: Update Cloudflare Workers (if using)

Edit `cloudflare/wrangler.toml`:

```toml
[env.production]
route = "your-domain.com/*"

[env.staging]
route = "staging.your-domain.com/*"
```

## 🔒 SSL Certificate Options

### Option 1: Let's Encrypt (Free)

```bash
# Install certbot
sudo apt-get install certbot

# Get certificate
sudo certbot certonly --standalone -d your-domain.com -d www.your-domain.com

# Copy certificates
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem nginx/ssl/key.pem
sudo chown $USER:$USER nginx/ssl/*
```

### Option 2: Cloudflare (Recommended)

1. Add your domain to Cloudflare
2. Enable SSL/TLS encryption mode to "Full (strict)"
3. Download certificates from Cloudflare dashboard
4. Place in `nginx/ssl/` directory

### Option 3: Commercial Certificate

Purchase from providers like:
- DigiCert
- GlobalSign
- Comodo

## 🌐 DNS Configuration

Configure your DNS records:

### A Records
```
your-domain.com     →  YOUR_SERVER_IP
www.your-domain.com →  YOUR_SERVER_IP
```

### CNAME Records (if using Cloudflare)
```
your-domain.com     →  your-domain.com.cdn.cloudflare.com
www.your-domain.com →  your-domain.com.cdn.cloudflare.com
```

## 🧪 Testing Your Configuration

### Test Domain Configuration

```bash
./scripts/test-domain.sh
```

This will test:
- ✅ DNS resolution
- ✅ SSL certificate presence
- ✅ Docker configuration validity

### Test SSL Certificate

```bash
# Test certificate validity
openssl x509 -in nginx/ssl/cert.pem -text -noout

# Test SSL connection
openssl s_client -connect your-domain.com:443 -servername your-domain.com
```

### Test PCO Integration

```bash
# Test PCO API connection
curl -H "Authorization: Bearer YOUR_PCO_ACCESS_TOKEN" \
     https://api.planningcenteronline.com/people/v2/me
```

## 🚀 Deployment

### 1. Pre-Deployment Checklist

- [ ] SSL certificates in place
- [ ] DNS configured correctly
- [ ] PCO credentials verified
- [ ] Environment variables set
- [ ] Docker and Docker Compose installed

### 2. Deploy to Production

```bash
./scripts/deploy-production.sh
```

### 3. Verify Deployment

- [ ] Application accessible at `https://your-domain.com`
- [ ] SSL certificate working
- [ ] PCO authentication working
- [ ] WebSocket connections working
- [ ] Admin panel accessible

## 🔧 Troubleshooting

### Common Issues

#### 1. SSL Certificate Issues
```bash
# Check certificate validity
openssl x509 -in nginx/ssl/cert.pem -text -noout

# Check private key
openssl rsa -in nginx/ssl/key.pem -check

# Test SSL connection
openssl s_client -connect your-domain.com:443
```

#### 2. DNS Issues
```bash
# Check DNS resolution
nslookup your-domain.com
dig your-domain.com

# Check propagation
https://www.whatsmydns.net/
```

#### 3. PCO Authentication Issues
```bash
# Test PCO API access
curl -H "Authorization: Bearer YOUR_TOKEN" \
     https://api.planningcenteronline.com/people/v2/me
```

#### 4. Docker Issues
```bash
# Check Docker configuration
docker-compose -f docker-compose.production.yml config

# Check container logs
docker-compose -f docker-compose.production.yml logs
```

### Debug Commands

```bash
# Check environment variables
grep -E "^(PCO_|CORS_|DOMAIN)" .env.production

# Check SSL certificate
ls -la nginx/ssl/

# Test Docker build
docker-compose -f docker-compose.production.yml build

# Check network connectivity
curl -I https://your-domain.com
```

## 📊 Monitoring

### Health Checks

- **Application**: `https://your-domain.com/health`
- **Grafana**: `http://localhost:3001`
- **Prometheus**: `http://localhost:9090`

### Logs

```bash
# Application logs
docker-compose -f docker-compose.production.yml logs -f

# Nginx logs
docker-compose -f docker-compose.production.yml logs nginx

# Database logs
docker-compose -f docker-compose.production.yml logs postgres
```

## 🔄 Maintenance

### SSL Certificate Renewal

For Let's Encrypt certificates:

```bash
# Renew certificates
sudo certbot renew

# Copy renewed certificates
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem nginx/ssl/key.pem

# Restart services
docker-compose -f docker-compose.production.yml restart nginx
```

### Regular Updates

```bash
# Update Docker images
docker-compose -f docker-compose.production.yml pull

# Rebuild and restart
docker-compose -f docker-compose.production.yml up -d --build
```

## 📞 Support

### Configuration Files

- **Environment**: `.env.production`
- **Docker Compose**: `docker-compose.production.yml`
- **Nginx**: `nginx/conf.d/your-domain.conf`
- **SSL**: `nginx/ssl/cert.pem`, `nginx/ssl/key.pem`

### Important URLs

- **Application**: `https://your-domain.com`
- **Admin Panel**: `https://your-domain.com/admin`
- **Billboard**: `https://your-domain.com/billboard`
- **Health Check**: `https://your-domain.com/health`

### Security Notes

- Keep `.env.production` secure and never commit to version control
- Regularly update SSL certificates
- Monitor logs for security issues
- Set up automated security updates
- Use strong, unique secrets for production

---

**Ready to deploy your PCO Arrivals Dashboard with your custom domain! 🚀** 