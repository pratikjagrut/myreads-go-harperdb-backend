FROM golang:alpine AS builder

ARG DB_HOST
ARG BASIC_AUTH_TOKEN
ARG HDB_ADMIN
ARG PASSWORD
ARG PORT=8000
ARG IMAGES_DIR=images

# Set necessary environmet variables needed for our image
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY . .
RUN go mod download

# Build the application
RUN go build -o myreads .


FROM alpine

ARG DB_HOST
ARG BASIC_AUTH_TOKEN
ARG HDB_ADMIN
ARG PASSWORD
ARG PORT=8000
ARG IMAGES_DIR=images

ENV HARPERDB_HOST=${DB_HOST}
ENV BASIC_AUTH_TOKEN=${BASIC_AUTH_TOKEN}
ENV HARPERDB_UNAME=${HDB_ADMIN}
ENV HARPERDB_PSWD=${PASSWORD}
ENV PORT=${PORT}
ENV IMAGES_DIR=${IMAGES_DIR}

WORKDIR /server

COPY --from=builder /build/myreads .

COPY  ./initiate_server.sh .

# Make initiate_server.sh executable
RUN chmod +x ./initiate_server.sh

# Install curl
RUN apk --no-cache add curl

RUN mkdir ${IMAGES_DIR}

# Expose necessary port
EXPOSE ${PORT}

# Command to run when starting the container
ENTRYPOINT [ "./initiate_server.sh" ]

# Override $1=false in docker run if DB is already initilized
CMD [ "false" ]
