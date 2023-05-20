# syntax = docker/dockerfile:experimental

FROM golang:1.20 AS build-stage

# Build the application from source
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY config ./config
COPY controller ./controller
COPY model ./model
COPY router ./router
COPY session ./session

RUN CGO_ENABLED=1 GOOS=linux go build -o /immich-goserver -a -ldflags '-linkmode external -extldflags "-static"' .

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage
SHELL ["busybox", "sh", "-c"]

WORKDIR /

COPY --from=build-stage /immich-goserver /immich-goserver

RUN --mount=from=busybox:uclibc,dst=/usr/ mkdir -p /var/lib/immich/database \
 && mkdir -p /var/lib/immich/upload \
 && mkdir -p /var/lib/immich/library \
 && mkdir -p /var/lib/immich/thumbnail \
 && mkdir -p /var/lib/immich/encoded-video

ENV GIN_MODE=release

EXPOSE 8080

ENTRYPOINT ["/immich-goserver"]
