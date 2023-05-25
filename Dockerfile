# syntax = docker/dockerfile:experimental

FROM golang:1.20.4-alpine3.18 AS build-stage

RUN apk add --no-cache gcc musl-dev exiftool

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY config ./config
COPY controller ./controller
COPY helper ./helper
COPY model ./model
COPY pipeline ./pipeline
COPY router ./router
COPY session ./session
COPY view ./view
COPY test ./test

RUN CGO_ENABLED=1 GOOS=linux go build -o /immich-goserver -a -ldflags '-linkmode external -extldflags "-static"' .

FROM build-stage AS run-test-stage
RUN GIN_MODE=release go test -v ./...

FROM alpine:3.18 AS build-release-stage

RUN apk add --no-cache exiftool

WORKDIR /

COPY --from=build-stage /immich-goserver /immich-goserver

RUN mkdir -p /var/lib/immich/database \
 && mkdir -p /var/lib/immich/upload \
 && mkdir -p /var/lib/immich/library \
 && mkdir -p /var/lib/immich/thumbnail \
 && mkdir -p /var/lib/immich/encoded-video

ENV GIN_MODE=release

EXPOSE 8080

ENTRYPOINT ["/immich-goserver"]
