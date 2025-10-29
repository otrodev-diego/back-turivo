#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================="
echo "  Restore Backup using DATABASE_URL"
echo "========================================="
echo ""

# Check if backup file exists
if [ -z "$1" ]; then
    echo -e "${RED}✗ Error: Please provide the backup file path${NC}"
    echo ""
    echo "Usage: $0 <backup_file.sql>"
    echo ""
    echo "Available backup files:"
    ls -lh turivo_backup_*.sql 2>/dev/null || echo "  No backup files found"
    exit 1
fi

BACKUP_FILE=$1

if [ ! -f "$BACKUP_FILE" ]; then
    echo -e "${RED}✗ Error: Backup file not found: $BACKUP_FILE${NC}"
    exit 1
fi

# Get DATABASE_URL from user
echo -e "${YELLOW}Please enter your Render DATABASE_URL:${NC}"
echo "Format: postgres://user:password@host:port/database"
echo ""
read -p "DATABASE_URL: " DATABASE_URL

if [ -z "$DATABASE_URL" ]; then
    echo -e "${RED}✗ Error: DATABASE_URL cannot be empty${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}Restoring backup...${NC}"
echo ""

# Restore using psql with DATABASE_URL
psql "$DATABASE_URL" < "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Backup restored successfully!${NC}"
    echo ""
else
    echo ""
    echo -e "${RED}✗ Error restoring backup${NC}"
    exit 1
fi
