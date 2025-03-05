FROM node:18-alpine AS base

# Install Python and build dependencies for node-gyp
# Also include git, in case npm packages require git to be fetched
RUN apk update && apk add --no-cache \
    python3 \
    make \
    g++ \
    git \
    && ln -sf python3 /usr/bin/python

# 1. Install dependencies only when needed
FROM base AS deps
# Check https://github.com/nodejs/docker-node/tree/b4117f9333da4138b03a546ec926ef50a31506c3#nodealpine to understand why libc6-compat might be needed.
RUN apk add --no-cache libc6-compat

WORKDIR /app

# Copy dependency definitions + rogue
COPY package.json package-lock.json . 
RUN npm ci

# 2. Rebuild the source code only when needed
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

ARG NEXT_PUBLIC_GOOGLE_CLIENT_ID=BAKED_NEXT_PUBLIC_GOOGLE_CLIENT_ID
ARG NEXT_PUBLIC_APP_HOST=BAKED_NEXT_PUBLIC_APP_HOST
ARG NEXT_PUBLIC_WS_HOST=BAKED_NEXT_PUBLIC_WS_HOST
ARG NEXT_PUBLIC_POSTHOG_HOST=BAKED_NEXT_PUBLIC_POSTHOG_HOST
ARG NEXT_PUBLIC_POSTHOG_KEY=BAKED_NEXT_PUBLIC_POSTHOG_KEY
ARG NEXT_PUBLIC_SEGMENT_WRITE_KEY=BAKED_NEXT_PUBLIC_SEGMENT_WRITE_KEY

RUN npm run build

# 3. Production image, copy all the files and run next
FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production

RUN addgroup -g 1001 -S nodejs
RUN adduser -S nextjs -u 1001

COPY --from=builder /app/public ./public

# Automatically leverage output traces to reduce image size
# https://nextjs.org/docs/advanced-features/output-file-tracing
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static
COPY --from=builder --chown=nextjs:nodejs /app/scripts ./scripts

USER nextjs

EXPOSE 3000

ENV PORT 3000
# set hostname to localhost
ENV HOSTNAME "0.0.0.0"

CMD ["/app/scripts/start.sh"]

