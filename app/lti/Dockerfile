FROM golang:1.18.1-alpine3.15 AS builder

# Set the working directory to /app
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the app
RUN go build -o librarian cmd/serve.go

# Start from a new image and copy the built executable into it
FROM alpine:3.15
COPY --from=builder /app/librarian /usr/local/bin/

# Set the default command to run when the container starts
CMD ["librarian"]