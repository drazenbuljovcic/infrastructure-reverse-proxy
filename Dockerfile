# Use the official Go image as the base image
FROM golang:1.17-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go service source code to the container's working directory
COPY . .

# Build the Go service inside the container
RUN go build -o go-service

# Expose the port that your Go service will listen on
EXPOSE 8080

# Set the command to run your Go service when the container starts
CMD ["./go-service"]

# docker build -t reverse-proxy .
# docker run -p 8080:8080 --env-file ../.env reverse-proxy

ENV SERVICE_HOSTNAME=http://localhost:3001