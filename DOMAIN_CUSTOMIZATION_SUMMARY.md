# 🌐 Domain Customization - Implementation Summary

## ✅ **What We've Implemented**

### 🚀 **Automated Domain Configuration Script**
- **File**: `scripts/customize-domain.sh`
- **Purpose**: Interactive script to configure your domain automatically
- **Features**:
  - ✅ Prompts for domain name and PCO credentials
  - ✅ Generates secure secrets automatically
  - ✅ Updates all configuration files
  - ✅ Creates SSL certificate instructions
  - ✅ Generates deployment checklist

### 📋 **Comprehensive Configuration Guide**
- **File**: `DOMAIN_CONFIGURATION_GUIDE.md`
- **Purpose**: Step-by-step manual configuration instructions
- **Features**:
  - ✅ Manual configuration steps
  - ✅ SSL certificate options (Let's Encrypt, Cloudflare, Commercial)
  - ✅ DNS configuration instructions
  - ✅ Testing procedures
  - ✅ Troubleshooting guide

### 🔧 **Domain-Specific Configuration Files**

#### 1. Environment Variables (`.env.production`)
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
```

#### 2. Docker Compose Configuration (`docker-compose.production.yml`)
```yaml
environment:
  - VITE_API_BASE_URL=${VITE_API_BASE_URL:-https://your-domain.com}
  - VITE_WS_BASE_URL=${VITE_WS_BASE_URL:-wss://your-domain.com}
```

#### 3. Nginx Configuration (`nginx/conf.d/your-domain.conf`)
- ✅ SSL/TLS termination
- ✅ HTTP to HTTPS redirect
- ✅ Security headers
- ✅ Rate limiting
- ✅ WebSocket proxy support

#### 4. Cloudflare Workers Configuration (`cloudflare/wrangler.toml`)
```toml
[env.production]
route = "your-domain.com/*"

[env.staging]
route = "staging.your-domain.com/*"
```

### 🔒 **SSL Certificate Management**
- **Directory**: `nginx/ssl/`
- **Files**: `cert.pem`, `key.pem`
- **Instructions**: `nginx/ssl/README.md`
- **Options**:
  - Let's Encrypt (free)
  - Cloudflare (recommended)
  - Commercial certificates

### 🧪 **Testing Tools**
- **File**: `scripts/test-domain.sh`
- **Tests**:
  - ✅ DNS resolution
  - ✅ SSL certificate presence
  - ✅ Docker configuration validity

### 📋 **Deployment Checklist**
- **File**: `DEPLOYMENT_CHECKLIST.md`
- **Features**:
  - ✅ Pre-deployment tasks
  - ✅ Deployment steps
  - ✅ Post-deployment verification
  - ✅ Security checklist

## 🚀 **How to Use**

### Option 1: Automated Configuration (Recommended)

```bash
# Run the automated configuration script
./scripts/customize-domain.sh
```

This will:
1. Prompt for your domain and PCO credentials
2. Generate secure secrets
3. Update all configuration files
4. Create SSL certificate instructions
5. Generate a deployment checklist

### Option 2: Manual Configuration

Follow the step-by-step guide in `DOMAIN_CONFIGURATION_GUIDE.md`

## 📋 **Required Information**

Before running the configuration, gather:

### Domain Information
- **Domain Name**: Your domain (e.g., `mychurch.org`)
- **SSL Certificates**: SSL certificate and private key files

### PCO API Credentials
- **PCO Client ID**: From your PCO application
- **PCO Client Secret**: From your PCO application  
- **PCO Access Token**: Your PCO access token
- **PCO Access Secret**: Your PCO access secret
- **PCO User ID**: Your PCO user ID for admin access

## 🔧 **Configuration Files Created/Updated**

### Core Configuration
- ✅ `.env.production` - Production environment variables
- ✅ `docker-compose.production.yml` - Docker orchestration
- ✅ `nginx/conf.d/your-domain.conf` - Nginx server configuration

### SSL Configuration
- ✅ `nginx/ssl/README.md` - SSL certificate instructions
- ✅ `nginx/ssl/cert.pem` - SSL certificate (you add)
- ✅ `nginx/ssl/key.pem` - SSL private key (you add)

### Testing & Documentation
- ✅ `scripts/test-domain.sh` - Domain testing script
- ✅ `DEPLOYMENT_CHECKLIST.md` - Deployment checklist
- ✅ `DOMAIN_CONFIGURATION_GUIDE.md` - Configuration guide

### Cloudflare Integration
- ✅ `cloudflare/wrangler.toml` - Cloudflare Workers configuration

## 🌐 **DNS Configuration Required**

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

## 🔒 **SSL Certificate Options**

### Option 1: Let's Encrypt (Free)
```bash
sudo certbot certonly --standalone -d your-domain.com -d www.your-domain.com
```

### Option 2: Cloudflare (Recommended)
1. Add domain to Cloudflare
2. Enable SSL/TLS encryption mode to "Full (strict)"
3. Download certificates from Cloudflare dashboard

### Option 3: Commercial Certificate
Purchase from providers like DigiCert, GlobalSign, Comodo

## 🧪 **Testing Your Configuration**

### Test Domain Configuration
```bash
./scripts/test-domain.sh
```

### Test SSL Certificate
```bash
openssl x509 -in nginx/ssl/cert.pem -text -noout
openssl s_client -connect your-domain.com:443
```

### Test PCO Integration
```bash
curl -H "Authorization: Bearer YOUR_PCO_ACCESS_TOKEN" \
     https://api.planningcenteronline.com/people/v2/me
```

## 🚀 **Deployment Steps**

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

## 📊 **Monitoring & Maintenance**

### Health Checks
- **Application**: `https://your-domain.com/health`
- **Grafana**: `http://localhost:3001`
- **Prometheus**: `http://localhost:9090`

### SSL Certificate Renewal
```bash
sudo certbot renew
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem nginx/ssl/key.pem
docker-compose -f docker-compose.production.yml restart nginx
```

## 🔧 **Troubleshooting**

### Common Issues
1. **SSL Certificate Issues**: Check certificate validity and permissions
2. **DNS Issues**: Verify DNS resolution and propagation
3. **PCO Authentication Issues**: Test PCO API access
4. **Docker Issues**: Check Docker configuration and logs

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

## 🎯 **Next Steps**

1. **Run Domain Configuration**:
   ```bash
   ./scripts/customize-domain.sh
   ```

2. **Add SSL Certificates**:
   - Place certificates in `nginx/ssl/`
   - Set permissions: `chmod 600 nginx/ssl/*`

3. **Configure DNS**:
   - Point your domain to your server IP
   - Verify DNS propagation

4. **Test Configuration**:
   ```bash
   ./scripts/test-domain.sh
   ```

5. **Deploy to Production**:
   ```bash
   ./scripts/deploy-production.sh
   ```

## 📞 **Support Information**

### Important URLs
- **Application**: `https://your-domain.com`
- **Admin Panel**: `https://your-domain.com/admin`
- **Billboard**: `https://your-domain.com/billboard`
- **Health Check**: `https://your-domain.com/health`

### Configuration Files
- **Environment**: `.env.production`
- **Docker Compose**: `docker-compose.production.yml`
- **Nginx**: `nginx/conf.d/your-domain.conf`
- **SSL**: `nginx/ssl/cert.pem`, `nginx/ssl/key.pem`

### Security Notes
- Keep `.env.production` secure and never commit to version control
- Regularly update SSL certificates
- Monitor logs for security issues
- Set up automated security updates
- Use strong, unique secrets for production

---

**Your PCO Arrivals Dashboard is ready for domain customization! 🚀**

Run `./scripts/customize-domain.sh` to get started with the automated configuration process. 