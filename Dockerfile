# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.19 AS build-stage

ENV GO111MODULE=on

ADD . /app
WORKDIR /app
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o /gcom-be

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /
COPY --from=build-stage /gcom-be /gcom-be
COPY --from=build-stage /app/db /db
EXPOSE 1323

ENTRYPOINT ["/gcom-be"]