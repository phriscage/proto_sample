###############################################################################
##  Name:   Dockerfile
##  Date:   2023-12-21
##  Developer:  Chris Page
##  Email:  phriscage@gmail.com
##  Purpose:   This Dockerfile contains the proto sample grpc server
################################################################################
# Using official golang alpine image
FROM --platform=linux/amd64 golang:1.21.3-alpine AS builder
# Set the file maintainer (your name - the file's author)
MAINTAINER Chris Page <phriscage@gmail.com>

# Set the GO build variables
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# Install git and ca-certificates (needed to be able to call HTTPS)
RUN apk add --no-cache ca-certificates git
#RUN apk add build-base

# Move to working directory /app
WORKDIR /app

# Copy the dependencies
COPY go.mod go.sum ./
# Download dependencies using go mod
RUN go mod download

# Copy the generated protobuf files into the appropriate package module directory
COPY gen/go/sample/v1alpha gen/go/sample/v1alpha

# Copy the code into the container
COPY server/ server

# Build the application's binary
# go build -tags musl --ldflags "-extldflags -static" -o main .
RUN go build -tags musl --ldflags "-extldflags -static" -o bin/samplesrv ./server


# Build a smaller image that will only contain the application's binary
FROM scratch

# Move to working directory /app
WORKDIR /app

# Copy certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# Copy application's binary
COPY --from=builder /app .

# Expose the port
EXPOSE 10000
# Command to run the application when starting the container
ENTRYPOINT ["/app/bin/samplesrv"]
CMD ["--host", "0.0.0.0"]
