# Use the official Go image as the base image
FROM golang:1.20.7-bookworm

# Set the working directory inside the container
WORKDIR /app

# Copy the Go service source code to the container's working directory
COPY . .

# Build the Go service inside the container
RUN RUN go build -o reverse-proxy

# Expose the port that your Go service will listen on
EXPOSE 8080

ENV PORT=8080
ENV SERVICE_HOSTNAME=http://ec2-13-48-148-19.eu-north-1.compute.amazonaws.com
ENV OTEL_SERVICE_NAME=infrastructure-reverse-proxy
ENV OTEL_API_HOST=http://ec2-16-171-166-86.eu-north-1.compute.amazonaws.com:9411

CMD ["chmod +x ./reverse-proxy", "./reverse-proxy"]

# docker build -t infrastructure-reverse-proxy . && docker run -d -p 8080:8080 -p 80:8080  --env-file .env infrastructure-reverse-proxy
