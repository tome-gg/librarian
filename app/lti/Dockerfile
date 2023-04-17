FROM arm64v8/golang:1.20.3 AS builder

# Set the working directory to /app
WORKDIR /usr/local/go/src

# Copy the source code into the container
COPY . .

# Build the app
RUN go build -o librarian cmd/serve.go

# Start from a new image and copy the built executable into it
FROM alpine:3.15
COPY --from=builder /usr/local/go/src/librarian /usr/local/bin/

# Set the default command to run when the container starts
CMD ["librarian"]