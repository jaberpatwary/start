#!/usr/bin/env bash
# ==============================================================================
# MI-Tech (StarTech Clone) — Production Deployment Script
# Usage: ./deploy.sh
# ==============================================================================

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}====================================================${NC}"
echo -e "${GREEN}🚀 Starting MI-Tech Production Deployment...${NC}"
echo -e "${GREEN}====================================================${NC}"

# Check Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

# Check Docker Compose
if ! docker compose version &> /dev/null; then
    echo -e "${RED}❌ Docker Compose plugin is not installed.${NC}"
    exit 1
fi

# Determine env file (.env.production or .env)
ENV_FILE=".env.production"
if [ ! -f "$ENV_FILE" ]; then
    if [ -f ".env" ]; then
        ENV_FILE=".env"
        echo -e "${YELLOW}⚠️  .env.production not found, falling back to .env${NC}"
    else
        echo -e "${RED}❌ Neither .env.production nor .env found!${NC}"
        echo -e "${YELLOW}👉 Run: cp .env.production.example .env.production and edit it.${NC}"
        exit 1
    fi
fi

echo -e "${GREEN}📦 Using environment file: ${ENV_FILE}${NC}"

# Optional git pull if in a git repository
if [ -d ".git" ]; then
    echo -e "${YELLOW}📥 Pulling latest updates from git repository...${NC}"
    git pull origin main || echo -e "${YELLOW}⚠️ Git pull skipped or failed, continuing with local code...${NC}"
fi

# Build containers
echo -e "${GREEN}🔨 Building production Docker containers...${NC}"
docker compose --env-file "$ENV_FILE" build --pull

# Start containers in detached mode
echo -e "${GREEN}🚀 Starting services with zero-downtime recreation...${NC}"
docker compose --env-file "$ENV_FILE" up -d --remove-orphans

# Clean up dangling images to free up disk space
echo -e "${YELLOW}🧹 Cleaning up dangling images...${NC}"
docker image prune -f

# Health check wait
echo -e "${YELLOW}⏳ Waiting for services to become healthy (15s)...${NC}"
sleep 15

# Verify backend health endpoint
echo -e "${GREEN}🩺 Checking system health...${NC}"
if curl -s -f http://localhost:8090/health > /dev/null 2>&1 || curl -s -f http://localhost/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ MI-Tech services are UP and HEALTHY!${NC}"
else
    echo -e "${YELLOW}⚠️ Health endpoint not responding immediately. Checking container status:${NC}"
fi

docker compose --env-file "$ENV_FILE" ps

echo -e "${GREEN}====================================================${NC}"
echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
echo -e "${GREEN}====================================================${NC}"
