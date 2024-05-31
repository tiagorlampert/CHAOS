# BUILD STAGE
FROM golang:1.22.3-alpine AS build

ARG APP_VERSION=dev
ARG CGO=1
ENV CGO_ENABLED=${CGO}
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on

# required by go-sqlite3
RUN apk add --update gcc g++ musl-dev

WORKDIR /build
COPY . .
RUN go build -v -a -tags 'netgo' -ldflags '-w -X 'main.Version=${APP_VERSION}' -linkmode external -extldflags -static' -o chaos cmd/chaos/*

# FINAL STAGE
FROM golang:1.22.3

MAINTAINER tiagorlampert@gmail.com

ENV GIN_MODE=release

WORKDIR /
COPY --from=build /build/chaos /
COPY ./web /web
COPY ./client /client

EXPOSE 8080
ENTRYPOINT ["/chaos"]
