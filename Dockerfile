# Build stage
FROM golang:1.25-alpine AS build

# Set the working directory
WORKDIR /build

# Install git
RUN --mount=type=cache,target=/var/cache/apk \
    apk add git

# Environment variables for Go
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Build the server
# go build automatically download required module dependencies to /go/pkg/mod
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,target=. \
    go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=$(git rev-parse HEAD) -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o /bin/mcp2grule main.go

# Make a stage to run the app
FROM alpine:3.22

# Set the working directory
WORKDIR /server

# Copy the binary from the build stage
COPY --from=build /bin/mcp2grule .

# Expose ports for mcp server and pprof
EXPOSE 9000 9001

# Set the entrypoint to the server binary
ENTRYPOINT ["/server/mcp2grule"]

# Default arguments for ENTRYPOINT
CMD ["server"]
