FROM golang:alpine

ARG DB_HOST
ARG BASIC_AUTH_TOKEN
ARG HDB_ADMIN
ARG PASSWORD

# Set necessary environmet variables needed for our image
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
ENV HARPERDB_HOST=${DB_HOST}
ENV BASIC_AUTH_TOKEN=${BASIC_AUTH_TOKEN}
ENV HARPERDB_UNAME=${HDB_ADMIN}
ENV HARPERDB_PSWD=${PASSWORD}

# Move to working directory /build
WORKDIR /server

# Copy and download dependency using go mod
ADD . .
RUN go mod download

# Make initiate_server.sh executable
RUN chmod +x ./initiate_server.sh

# Build the application
RUN go build -o myreads .

# Install curl
RUN apk --no-cache add curl

# Export necessary port
EXPOSE 8000

# Command to run when starting the container
ENTRYPOINT [ "./initiate_server.sh" ]