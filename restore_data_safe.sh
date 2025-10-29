#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

DATABASE_URL="postgresql://turivo_user:QU8FSmMlcxs3hd44S60HoLY2PHB18Zxe@dpg-d39o8eadbo4c73cj59p0-a.oregon-postgres.render.com:5432/turivo"
DATA_FILE="turivo_data_only.sql"

echo "========================================="
echo "  Safe Data Restore to Render"
echo "========================================="
echo ""
echo -e "${YELLOW}⚠️  WARNING: This will DELETE all existing data!${NC}"
echo ""
read -p "Are you sure you want to continue? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo -e "${RED}Restore cancelled${NC}"
    exit 0
fi

echo ""
echo -e "${YELLOW}Step 1: Truncating all tables...${NC}"

# Truncate all tables except schema_migrations
psql "$DATABASE_URL" << SQL
DO \$\$ 
DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename != 'schema_migrations') 
    LOOP
        EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' RESTART IDENTITY CASCADE';
    END LOOP;
END \$\$;
SQL

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Tables truncated${NC}"
else
    echo -e "${RED}✗ Error truncating tables${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}Step 2: Restoring data...${NC}"

psql "$DATABASE_URL" < "$DATA_FILE"

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Data restored successfully!${NC}"
else
    echo ""
    echo -e "${RED}✗ Error restoring data${NC}"
    exit 1
fi

echo ""
echo "========================================="
