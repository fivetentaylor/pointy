#!/bin/sh
# start.sh

# server.js is created by next build from the standalone output
# https://nextjs.org/docs/pages/api-reference/next-config-js/output

# Replace runtime env vars and start next server
sh scripts/replace-variables.sh && 
node server.js
