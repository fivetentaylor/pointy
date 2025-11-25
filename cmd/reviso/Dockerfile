FROM --platform=linux/amd64 alpine:latest as builder

RUN apk add --no-cache \
    fuse-overlayfs \
    bash \
    build-base \
    ca-certificates \
    cmake \
    curl \
    docker \
    docker-cli \
    docker-compose \
    git \
    glib-dev \
    gpgme-dev \
    iptables \
    libjpeg-turbo-dev \
    libpng-dev \
    npm \
    nss-dev \
    openjpeg-dev \
    pkgconfig \
    poppler-utils \
    tiff-dev \
    chromium \ 
    chromium-chromedriver \
    unrtf \
    && update-ca-certificates

# Install Go
ENV GO_VERSION=1.25.3
RUN curl -OL https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

WORKDIR /app

# Copying go.mod and go.sum first to leverage Docker cache
COPY go.mod .
COPY go.sum .

# Install dependencies - leverages cache if no changes to mod/sum
RUN go mod download

# Copy the rest of the source code
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY rogue/ rogue/
COPY gqlgen.yml .
COPY tools.go .

# Setup env vars
ARG NODE_ENV
ARG WEB_HOST
ARG APP_HOST
ARG WS_HOST
ARG SEGMENT_KEY
ARG PUBLIC_POSTHOG_KEY
ARG PUBLIC_POSTHOG_HOST
ARG IMAGE_TAG

# Build the wasm
RUN GOOS=js GOARCH=wasm go build -ldflags "-s -w -X main.ImageTag=$IMAGE_TAG" -o pkg/assets/static/rogueV3.wasm rogue/v3/wasm/main.go
RUN cp pkg/assets/static/rogueV3.wasm pkg/admin/static/rogue.wasm || echo "Warning: rogue.wasm copy failed"
RUN cp pkg/assets/static/rogueV3.wasm pkg/admin/static/rogueV3.wasm

# Install npm dependencies including dev dependencies
RUN (cd pkg/assets/src && npm install --include=dev) \
    || { echo 'npm install failed'; exit 1; }

# Build the frontend (so it can be embedded in the server binary)
ENV NODE_ENV=$NODE_ENV
ENV WEB_HOST=$WEB_HOST
ENV APP_HOST=$APP_HOST
ENV WS_HOST=$WS_HOST
ENV SEGMENT_KEY=$SEGMENT_KEY
ENV PUBLIC_POSTHOG_KEY=$PUBLIC_POSTHOG_KEY
ENV PUBLIC_POSTHOG_HOST=$PUBLIC_POSTHOG_HOST
ENV IMAGE_TAG=$IMAGE_TAG
RUN CGO_ENABLED=0 go run ./cmd/reviso/main.go --build

# Print out the dist folder
RUN ls -al /app/pkg/assets/src/dist

# Remove this if you want to run it locally
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Package up the server
RUN go build -ldflags "-s -w -X main.ImageTag=$IMAGE_TAG" -o reviso ./cmd/reviso/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates, chromium, chromium-chromedriver, nss, and poppler-utils
RUN apk --no-cache add ca-certificates chromium chromium-chromedriver nss curl poppler-utils unrtf

# Create a directory for custom fonts
RUN mkdir -p /usr/share/fonts/noto

# Copy the Noto Color Emoji font into the Docker image
COPY fonts/NotoColorEmoji.ttf /usr/share/fonts/noto

# Update font cache
RUN fc-cache -f -v

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/reviso .

# Run the server binary
CMD ["./reviso"]
