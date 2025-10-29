#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================="
echo "  Restore Backup to Render PostgreSQL"
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

echo -e "${YELLOW}Please enter your Render PostgreSQL connection details:${NC}"
echo ""

# Get connection details from user
read -p "Render PostgreSQL Host: " RENDER_HOST
read -p "Render PostgreSQL Port [5432]: " RENDER_PORT
read -p "Render PostgreSQL User: " RENDER_USER
read -p "Render PostgreSQL Database: " RENDER_DB
read -sp "Render PostgreSQL Password: " RENDER_PASSWORD
echo ""

# Set defaults
RENDER_PORT=${RENDER_PORT:-5432}

echo ""
echo -e "${YELLOW}Restoring backup to Render...${NC}"
echo ""

# Restore database
PGPASSWORD=$RENDER_PASSWORD psql \
  -h "$RENDER_HOST" \
  -p "$RENDER_PORT" \
  -U "$RENDER_USER" \
  -d "$RENDER_DB" \
  < "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Backup restored successfully to Render!${NC}"
    echo ""
else
    echo ""
    echo -e "${RED}✗ Error restoring backup${NC}"
    exit 1
fi
