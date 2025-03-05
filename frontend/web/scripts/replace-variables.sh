#!/bin/sh
# replace-variables.sh

# Define a list of environment variables to check and replace
# In sh, we'll handle this as a space-separated string since sh does not support arrays
VARIABLES="NEXT_PUBLIC_GOOGLE_CLIENT_ID NEXT_PUBLIC_APP_HOST NEXT_PUBLIC_WS_HOST NEXT_PUBLIC_POSTHOG_HOST NEXT_PUBLIC_POSTHOG_KEY NEXT_PUBLIC_SEGMENT_WRITE_KEY"

# Check if each variable is set
for VAR in $VARIABLES; do
    echo "Checking if $VAR is set..."
    eval VAR_VALUE=\$$VAR
    if [ -z "$VAR_VALUE" ]; then
        echo "$VAR is not set. Please set it and rerun the script."
        exit 1
    fi
done

# Find and replace BAKED values with real values
find /app/public /app/.next -type f -name "*.js" | while read file; do
    for VAR in $VARIABLES; do
        eval VAR_VALUE=\$$VAR
        sed -i "s|BAKED_$VAR|$VAR_VALUE|g" "$file"
    done
done
