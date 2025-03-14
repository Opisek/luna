# Base
FROM oven/bun:alpine AS base
WORKDIR /app

# Install bun production dependencies
FROM base AS dependencies
COPY package.json bun.lock ./
RUN bun install --production

# Install bun development dependencies
FROM dependencies AS devdependencies
RUN bun install

# Build the frontend
FROM base AS build
WORKDIR /app
RUN apk add make
COPY Makefile .env .npmrc svelte.config.js tsconfig.json VERSION vite.config.ts package.json ./
COPY script script
COPY static static
COPY src src
COPY --from=devdependencies /app/node_modules ./node_modules
RUN make build

# Run the frontend
FROM base AS run
WORKDIR /app
COPY package.json .
COPY --from=dependencies /app/node_modules ./node_modules
COPY --from=build /app/build ./build
USER bun:bun
ENTRYPOINT [ "bun", "./build/index.js" ]