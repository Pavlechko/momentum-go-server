FROM golang:bullseye AS builder

# Build Delve for debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Create and change to the app directory.
WORKDIR /app
ENV CGO_ENABLED=0

# Retrieve application dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
# COPY . ./
COPY . .

# Build the binary.
RUN go build -o fooapp ./cmd/momentum

# Download certificates
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates 

# Use the official Debian slim image for a lean production container.
FROM busybox:glibc

EXPOSE 8080 40000

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/fooapp /app/fooapp 
# COPY --from=builder /app/ /app

COPY --from=builder /go/bin/dlv /dlv

COPY --from=builder /etc/ssl /etc/ssl

# Run dlv as pass fooapp as parameter
CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/fooapp"]
# ENTRYPOINT ["/bin/sh"]
