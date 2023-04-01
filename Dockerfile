# Use the official Golang image as a base
FROM golang:1.17

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go server
RUN go build -o wifi-control-server

# Expose the server's port
EXPOSE 8080

# Set environment variables
ENV ROUTER_URL=""
ENV USERNAME=""
ENV PASSWORD=""

# Run the server
CMD ["/app/wifi-control-server"]
