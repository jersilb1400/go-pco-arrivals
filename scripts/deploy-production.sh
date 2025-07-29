#!/bin/bash

# PCO Arrivals Dashboard Production Deployment Script
# This script deploys the application to production

set -e

# Configuration
ENVIRONMENT="production"
COMPOSE_FILE="docker-compose.production.yml"
BACKUP_BEFORE_DEPLOY=true

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

# Check if running as root
if [ "$EUID" -eq 0 ]; then
    error_exit "Please do not run this script as root"
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    error_exit "Docker is not running"
fi

# Check if Docker Compose is available
if ! command -v docker-compose >/dev/null 2>&1; then
    error_exit "Docker Compose is not installed"
fi

# Check if required files exist
if [ ! -f "$COMPOSE_FILE" ]; then
    error_exit "Docker Compose file not found: $COMPOSE_FILE"
fi

if [ ! -f ".env.production" ]; then
    error_exit "Production environment file not found: .env.production"
fi

log "${BLUE}=== PCO Arrivals Dashboard Production Deployment ===${NC}"
log "Environment: $ENVIRONMENT"
log "Compose file: $COMPOSE_FILE"

# Pre-deployment checks
log "${YELLOW}Running pre-deployment checks...${NC}"

# Check disk space
DISK_USAGE=$(df / | tail -1 | awk '{print $5}' | sed 's/%//')
if [ "$DISK_USAGE" -gt 90 ]; then
    error_exit "Disk usage is too high: ${DISK_USAGE}%"
fi

# Check memory
MEMORY_AVAILABLE=$(free -m | awk 'NR==2{printf "%.0f", $7*100/$2}')
if [ "$MEMORY_AVAILABLE" -lt 20 ]; then
    error_exit "Available memory is too low: ${MEMORY_AVAILABLE}%"
fi

log "${GREEN}Pre-deployment checks passed${NC}"

# Create backup before deployment
if [ "$BACKUP_BEFORE_DEPLOY" = true ]; then
    log "${YELLOW}Creating backup before deployment...${NC}"
    ./scripts/backup.sh || log "${YELLOW}Warning: Backup failed, continuing with deployment${NC}"
fi

# Stop existing services
log "${YELLOW}Stopping existing services...${NC}"
docker-compose -f "$COMPOSE_FILE" down --remove-orphans || true

# Pull latest images
log "${YELLOW}Pulling latest images...${NC}"
docker-compose -f "$COMPOSE_FILE" pull

# Build images
log "${YELLOW}Building images...${NC}"
docker-compose -f "$COMPOSE_FILE" build --no-cache

# Start services
log "${YELLOW}Starting services...${NC}"
docker-compose -f "$COMPOSE_FILE" up -d

# Wait for services to be healthy
log "${YELLOW}Waiting for services to be healthy...${NC}"
sleep 30

# Check service health
log "${YELLOW}Checking service health...${NC}"
if ! docker-compose -f "$COMPOSE_FILE" ps | grep -q "Up"; then
    error_exit "Some services failed to start"
fi

# Run health checks
log "${YELLOW}Running health checks...${NC}"
HEALTH_CHECK_RETRIES=10
HEALTH_CHECK_INTERVAL=10

for i in $(seq 1 $HEALTH_CHECK_RETRIES); do
    if curl -f http://localhost/health > /dev/null 2>&1; then
        log "${GREEN}Health check passed${NC}"
        break
    else
        if [ $i -eq $HEALTH_CHECK_RETRIES ]; then
            error_exit "Health check failed after $HEALTH_CHECK_RETRIES attempts"
        fi
        log "${YELLOW}Health check attempt $i/$HEALTH_CHECK_RETRIES failed, retrying in ${HEALTH_CHECK_INTERVAL}s...${NC}"
        sleep $HEALTH_CHECK_INTERVAL
    fi
done

# Show deployment status
log "${GREEN}=== Deployment Status ===${NC}"
docker-compose -f "$COMPOSE_FILE" ps

# Show service logs
log "${BLUE}=== Recent Service Logs ===${NC}"
docker-compose -f "$COMPOSE_FILE" logs --tail=20

# Show resource usage
log "${BLUE}=== Resource Usage ===${NC}"
docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}"

log "${GREEN}=== Production Deployment Completed Successfully! ===${NC}"
log "Application is now running at: https://your-domain.com"
log "Monitoring dashboard: http://localhost:3001 (Grafana)"
log "Metrics endpoint: http://localhost:9090 (Prometheus)" 