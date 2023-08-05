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

ENV PORT=8080
ENV SERVICE_HOSTNAME=http://ec2-13-48-148-19.eu-north-1.compute.amazonaws.com
ENV OTEL_SERVICE_NAME=infrastructure-reverse-proxy
ENV OTEL_API_HOST=http://ec2-16-171-166-86.eu-north-1.compute.amazonaws.com:9411

# Set the command to run your Go service when the container starts
CMD ["./go-service"]

# docker build -t reverse-proxy .
# docker run -p 8080:8080 --env-file ../.env reverse-proxy
