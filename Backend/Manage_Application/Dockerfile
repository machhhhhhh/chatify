# Stage 1: Build Stage
FROM golang:alpine AS builder

# Set the working directory in the container
WORKDIR /app

# Copy the go.mod and go.sum files first and download the dependencies.
# This is done first to leverage Docker cache layers, as these files 
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the entire project 
COPY . .

# Build the application
RUN go build -o main .

# Stage 2: Run-Time Stage
FROM alpine

WORKDIR /app

# Copy only the build file
COPY --from=builder /app/main /app

# Declare User and Group
ARG user=swadmin
ARG group=swgroup

# Create new user and add to group
RUN addgroup -g 1000 $group
RUN adduser -D -G $group $user

RUN apk update

# Set ownership (permission) of the application directory
RUN chown -R $user:$group /app

# Switch to the new user
USER $user

# Make port 3000 available to the world outside this container
EXPOSE 3000

# Run the executable
CMD ["./main"]