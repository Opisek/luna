# Base
FROM oven/bun:alpine AS base
WORKDIR /app

# Install bun production dependencies
FROM base AS dependencies
COPY package.json bun.lock ./
# TODO: figure out why we can't use the --production flag
#RUN bun install --production
RUN bun install

# Install bun development dependencies
FROM dependencies AS devDependencies
RUN bun install

# Build the frontend
FROM base AS build
WORKDIR /app
RUN apk add make
COPY Makefile .env .npmrc svelte.config.js tsconfig.json VERSION vite.config.ts package.json ./
COPY script script
COPY static static
COPY src src
COPY --from=devDependencies /app/node_modules ./node_modules
RUN make build

# Run the frontend
FROM base AS run
WORKDIR /app
COPY package.json .
COPY --from=dependencies /app/node_modules ./node_modules
COPY --from=build /app/build ./build
USER bun:bun
ENTRYPOINT [ "bun", "./build/index.js" ]