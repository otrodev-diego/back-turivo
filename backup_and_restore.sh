#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================="
echo "  Turivo Database Backup & Restore"
echo "========================================="
echo ""

# Local database configuration
LOCAL_HOST="localhost"
LOCAL_PORT="5432"
LOCAL_DB="turivo"
LOCAL_USER="postgres"
LOCAL_PASSWORD="postgres"

# Backup file
BACKUP_FILE="turivo_backup_$(date +%Y%m%d_%H%M%S).sql"
TIMESTAMP=$(date +%Y-%m-%d_%H-%M-%S)

echo -e "${YELLOW}Step 1: Creating backup from local database...${NC}"
PGPASSWORD=$LOCAL_PASSWORD pg_dump -h $LOCAL_HOST -p $LOCAL_PORT -U $LOCAL_USER -d $LOCAL_DB \
  --no-owner --no-acl --verbose \
  > "$BACKUP_FILE" 2>&1

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Backup created successfully: $BACKUP_FILE${NC}"
    echo ""
    
    # Show file size
    FILESIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    echo -e "${GREEN}  File size: $FILESIZE${NC}"
    echo ""
else
    echo -e "${RED}✗ Error creating backup${NC}"
    exit 1
fi

echo "========================================="
echo "  Next Steps:"
echo "========================================="
echo ""
echo "1. Get your Render PostgreSQL connection string:"
echo "   Render Dashboard > PostgreSQL > Connect"
echo ""
echo "2. Run restore with:"
echo "   ./restore_to_render.sh"
echo ""
echo "   Or manually restore using:"
echo "   PGPASSWORD=your_password psql -h your-host -p 5432 -U your_user -d your_db < $BACKUP_FILE"
echo ""
echo "========================================="
