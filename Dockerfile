# Use the official Go image as the base image
FROM go-client-app:latest as builder

# Set the working directory inside the container
WORKDIR /app

# Enable CGO for plugin compilation
ENV CGO_ENABLED=1

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code to the container
COPY . .

# Create a directory for the built plugins
RUN mkdir -p ./plugins

# Build the plugins using a bash loop
RUN for dir in yfinance_news yfinance python chrome example; do \
        cd "$dir" && \
        for file in *.go; do \
            basename=$(basename "$file" .go); \
            echo "Building plugin for $file"; \
            GOARCH=amd64 GOOS=linux go build -gcflags="all=-N -l" -buildmode=plugin -o /app/plugins/"$basename".so "$file"; \
        done && \
        cd ..; \
    done

# Final stage: copy the plugins to the final image
FROM go-client-app:latest
WORKDIR /app
COPY --from=builder /app/plugins /app/plugins

# Set the entrypoint to show built plugins (for testing purposes)
CMD ls /app/plugins
