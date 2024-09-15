# Install node production dependencies
FROM node:20-alpine as dependencies
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install --production

# Install node development dependencies
FROM dependencies as devDependencies
RUN npm install

# Build the frontend
FROM node:20-alpine as build
WORKDIR /app
RUN apk add make
COPY Makefile .env .npmrc svelte.config.js tsconfig.json VERSION vite.config.ts package.json .
COPY script script
COPY static static
COPY src src
COPY --from=devDependencies /app/node_modules ./node_modules
RUN make build

# Run the frontend
FROM node:20-alpine as run
WORKDIR /app
COPY package.json .
COPY --from=dependencies /app/node_modules ./node_modules
COPY --from=build /app/build ./build
USER node:node
ENTRYPOINT [ "node", "build/index.js" ]

