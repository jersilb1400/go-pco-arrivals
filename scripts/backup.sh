#!/bin/sh

# PCO Arrivals Dashboard Backup Script
# This script creates backups of the database and configuration files

set -e

# Configuration
BACKUP_DIR="/app/backups"
DATA_DIR="/app/data"
RETENTION_DAYS=${BACKUP_RETENTION_DAYS:-30}
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_NAME="pco_arrivals_backup_${TIMESTAMP}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1"
}

# Error handling
error_exit() {
    log "${RED}ERROR: $1${NC}" >&2
    exit 1
}

# Check if backup directory exists
if [ ! -d "$BACKUP_DIR" ]; then
    log "${YELLOW}Creating backup directory...${NC}"
    mkdir -p "$BACKUP_DIR"
fi

# Create backup
log "${GREEN}Starting backup: $BACKUP_NAME${NC}"

# Create temporary directory for backup
TEMP_DIR="/tmp/$BACKUP_NAME"
mkdir -p "$TEMP_DIR"

# Copy database files
if [ -f "$DATA_DIR/pco_billboard.db" ]; then
    log "Backing up database..."
    cp "$DATA_DIR/pco_billboard.db" "$TEMP_DIR/"
    
    # Create SQL dump for better portability
    if command -v sqlite3 >/dev/null 2>&1; then
        log "Creating SQL dump..."
        sqlite3 "$DATA_DIR/pco_billboard.db" ".dump" > "$TEMP_DIR/pco_billboard.sql"
    fi
else
    log "${YELLOW}Warning: Database file not found${NC}"
fi

# Copy configuration files (if they exist)
if [ -f "/app/.env.production" ]; then
    log "Backing up configuration..."
    cp /app/.env.production "$TEMP_DIR/" 2>/dev/null || true
fi

# Create backup archive
log "Creating backup archive..."
cd /tmp
tar -czf "$BACKUP_DIR/$BACKUP_NAME.tar.gz" "$BACKUP_NAME"
rm -rf "$TEMP_DIR"

# Verify backup
if [ -f "$BACKUP_DIR/$BACKUP_NAME.tar.gz" ]; then
    BACKUP_SIZE=$(du -h "$BACKUP_DIR/$BACKUP_NAME.tar.gz" | cut -f1)
    log "${GREEN}Backup completed successfully: $BACKUP_NAME.tar.gz (${BACKUP_SIZE})${NC}"
else
    error_exit "Backup file was not created"
fi

# Clean up old backups
log "Cleaning up old backups (older than $RETENTION_DAYS days)..."
find "$BACKUP_DIR" -name "pco_arrivals_backup_*.tar.gz" -type f -mtime +$RETENTION_DAYS -delete

# List remaining backups
BACKUP_COUNT=$(find "$BACKUP_DIR" -name "pco_arrivals_backup_*.tar.gz" | wc -l)
log "${GREEN}Backup cleanup completed. $BACKUP_COUNT backups remaining.${NC}"

# Show backup summary
log "${GREEN}=== Backup Summary ===${NC}"
log "Backup file: $BACKUP_NAME.tar.gz"
log "Location: $BACKUP_DIR"
log "Size: $BACKUP_SIZE"
log "Retention: $RETENTION_DAYS days"
log "Total backups: $BACKUP_COUNT"

log "${GREEN}Backup process completed successfully!${NC}" 