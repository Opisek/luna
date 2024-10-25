# Base image
FROM golang:1.23-alpine AS base
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

# Run the binary
FROM alpine:latest
WORKDIR /app
RUN apk add tzdata
COPY --from=build /app/src/luna-backend .
ENTRYPOINT [ "/app/luna-backend" ]