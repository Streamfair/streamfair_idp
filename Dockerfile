# Build Stage
FROM golang:1.22.0-alpine3.19 AS build
WORKDIR /streamfair_identity_provider
COPY . .
RUN go build -o identity_provider main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

# Run Stage
FROM alpine:3.19
WORKDIR /streamfair_identity_provider

# Copy the binary from the build stage
COPY --from=build /streamfair_identity_provider/identity_provider .
# Copy the downloaded migration binary from the build stage
COPY --from=build /streamfair_identity_provider/migrate ./migrate

COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8081
EXPOSE 9091

CMD [ "/streamfair_identity_provider/identity_provider" ]
ENTRYPOINT [ "/streamfair_identity_provider/start.sh" ]
