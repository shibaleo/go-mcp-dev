#!/bin/bash

# Load environment variables from .env.production
ENV_FILE=".env.production"

if [ ! -f "$ENV_FILE" ]; then
  echo "Error: $ENV_FILE not found"
  exit 1
fi

# Build --env flags from .env.production
ENV_FLAGS=""
while IFS='=' read -r key value; do
  # Skip empty lines and comments
  [[ -z "$key" || "$key" =~ ^# ]] && continue
  # Remove any surrounding quotes from value
  value="${value%\"}"
  value="${value#\"}"
  ENV_FLAGS="$ENV_FLAGS --env $key=$value"
done < "$ENV_FILE"

APP_NAME="go-mcp-dev"

# Create app (ignore error if already exists)
koyeb app create "$APP_NAME" 2>/dev/null || true

# Create service with all environment variables
koyeb service create "$APP_NAME" \
  --app "$APP_NAME" \
  --git github.com/shibaleo/go-mcp-dev \
  --git-branch main \
  --git-builder docker \
  --instance-type free \
  --ports 8080:http \
  --routes /:8080 \
  $ENV_FLAGS

echo "Deployment started. Check status with: koyeb service get $APP_NAME"
