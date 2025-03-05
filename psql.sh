#!/bin/bash

# Read DATABASE_URL from .env file
source .env

# Parse DATABASE_URL using a regular expression
if [[ $DATABASE_URL =~ postgresql://([^:]+):([^@]+)@([^:]+):([^/]+)/([^?]+).*sslmode=([^&]+) ]]; then
    USER=${BASH_REMATCH[1]}
    PASSWORD=${BASH_REMATCH[2]}
    HOST=${BASH_REMATCH[3]}
    PORT=${BASH_REMATCH[4]}
    DB_NAME=${BASH_REMATCH[5]}
    SSL_MODE=${BASH_REMATCH[6]}

    # Connect to the database using psql
    PGPASSWORD=$PASSWORD psql -h $HOST -p $PORT -U $USER -d $DB_NAME
else
    echo "Failed to parse DATABASE_URL"
fi

