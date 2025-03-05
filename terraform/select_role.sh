#!/bin/bash

# Get the directory containing the script
BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Help flag check and usage message
usage() {
    echo "Usage: $0 [role-arn]"
    echo
    echo "Arguments:"
    echo "  role-arn       AWS Role ARN to assume (defaults to DEFAULT_ROLE_ARN from .env)"
    echo
    echo "Options:"
    echo "  -h, --help    Show this help message"
    echo
    echo "Example:"
    echo "  $0 arn:aws:iam::123456789012:role/MyRole"
    echo "  $0  # Uses base credentials from terraform/.env"
    exit 1
}

# Check for help flags or invalid arguments
if [ "$1" = "-h" ] || [ "$1" = "--help" ] || [ "$#" -gt 1 ]; then
    usage
fi

# Assign arguments to variables
DEFAULT_ROLE_ARN="arn:aws:iam::533267310428:role/TerraformRole"
ROLE_ARN="${1:-$DEFAULT_ROLE_ARN}"

USERNAME=$(whoami)
TIMESTAMP=$(date '+%Y_%m%d_%H%M')
SESSION_NAME="awscli_${USERNAME}_${TIMESTAMP}"

# Check for .env file
if [ ! -f "${BASE_DIR}/.env" ]; then
    echo "Error: .env file not found in current directory"
    exit 1
fi

# Source the .env file
export $(grep -v '#' "${BASE_DIR}/.env" | xargs)

# Validate required environment variables
if [ -z "$AWS_ACCESS_KEY_ID" ] || [ -z "$AWS_SECRET_ACCESS_KEY" ]; then
    echo "Error: AWS credentials not found in ${BASE_DIR}/.env file"
    echo "Please ensure your .env file contains:"
    echo "AWS_ACCESS_KEY_ID=your_access_key"
    echo "AWS_SECRET_ACCESS_KEY=your_secret_key"
    exit 1
fi

# Unset existing session variables
unset AWS_SESSION_TOKEN

# Validate AWS CLI installation
if ! command -v aws &> /dev/null; then
    echo "Error: AWS CLI is not installed"
    exit 1
fi

# Get credentials for the specified role
CREDENTIALS=$(aws sts assume-role \
    --role-arn "$ROLE_ARN" \
    --role-session-name "$SESSION_NAME" \
    --output json 2>&1)

# Check if the assume-role command was successful
if [ $? -ne 0 ]; then
    echo "Error assuming role: $CREDENTIALS" >&2
    exit 1
fi

# Extract credentials from JSON response and print as export commands
echo "# Run these commands to switch to the new role:"
echo "export AWS_ACCESS_KEY_ID=$(echo "$CREDENTIALS" | grep -o '"AccessKeyId": "[^"]*' | cut -d'"' -f4)"
echo "export AWS_SECRET_ACCESS_KEY=$(echo "$CREDENTIALS" | grep -o '"SecretAccessKey": "[^"]*' | cut -d'"' -f4)"
echo "export AWS_SESSION_TOKEN=$(echo "$CREDENTIALS" | grep -o '"SessionToken": "[^"]*' | cut -d'"' -f4)"
echo
echo "# Or source them all at once using:"
echo "# eval \"\$($0 $1)\""
