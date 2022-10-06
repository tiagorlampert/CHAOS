# BUILD STAGE
FROM golang:1.18-alpine AS build

ARG APP_VERSION=dev
ARG CGO=1
ENV CGO_ENABLED=${CGO}
ENV GOOS=linux
ENV GO111MODULE=on

# gcc/g++ are required by sqlite driver
RUN apk update && apk add --no-cache gcc g++

WORKDIR /build
COPY . .
RUN go build -v -a -tags 'netgo' -ldflags '-w -X 'main.Version=${APP_VERSION}' -extldflags "-static"' -o chaos cmd/chaos/*

# FINAL STAGE
FROM golang:1.18.4

MAINTAINER tiagorlampert@gmail.com

ENV GIN_MODE=release

WORKDIR /
COPY --from=build /build/chaos /
COPY ./web /web
COPY ./client /client

EXPOSE 8080
ENTRYPOINT ["/chaos"]
