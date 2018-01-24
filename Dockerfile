FROM golang:1.9.2-alpine AS build

# Build args
ARG PROJECT
ENV PROJECT_SRC=/go/src/${PROJECT}

# Install tools required to build the project
# We will need to run `docker build --no-cache .` to update those dependencies
RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

# Gopkg.toml and Gopkg.lock lists project dependencies
# These layers will only be re-built when Gopkg files are updated
COPY Gopkg.lock Gopkg.toml ${PROJECT_SRC}/
WORKDIR ${PROJECT_SRC}

# Install library dependencies
RUN dep ensure -vendor-only

# Copy all project and build it
# This layer will be rebuilt when ever a file has changed in the project directory
COPY . ${PROJECT_SRC}/
RUN go build -o /project/server

# Copy config file to output directory
COPY ./config.yml /project/

# This results in a small single layer image
FROM alpine:latest

# Move output directory from the previous build to the new one
COPY --from=build /project /project

WORKDIR /project
ENTRYPOINT ["./server"]
