#!/bin/sh

if [ "$DEBUG" = "true" ]; then
    echo "Starting application in debug mode with Delve..."
    dlv exec /app/main --headless --listen=:40000 --api-version=2 --accept-multiclient
else
    echo "Starting application in development mode with Air..."
    air -c .air.toml
fi