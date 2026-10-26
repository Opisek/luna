# Base image
FROM golang:1.27-alpine AS base
WORKDIR /app
RUN apk add make

# Install go dependencies
FROM base AS dependencies
COPY src/go.mod src/go.sum ./
RUN go mod download

# Build the binary
FROM dependencies AS build
COPY Makefile .
COPY VERSION .
COPY src src
RUN make build

# Install runtime dependencies
FROM alpine:latest AS runtime
RUN echo 'http://dl-cdn.alpinelinux.org/alpine/edge/main' > /etc/apk/repositories
RUN apk update
RUN apk upgrade
RUN apk add postgresql18-client

# Run the binary
FROM runtime AS final
WORKDIR /app
RUN apk add --no-cache tzdata
COPY --from=build /app/src/luna-backend .
EXPOSE 3000
VOLUME [ "/data" ]
HEALTHCHECK --interval=15s --timeout=3s --retries=2 CMD wget -O /dev/null http://localhost:${PORT:-3000}/api/health
ENTRYPOINT [ "/app/luna-backend" ]
