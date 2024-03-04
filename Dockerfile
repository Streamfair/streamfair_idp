# Build Stage
FROM golang:1.22.0-alpine3.19 AS build
WORKDIR /streamfair_identity_provider
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o identity_provider main.go

# Run Stage
FROM alpine:3.19
WORKDIR /streamfair_identity_provider

# Copy the binary from the build stage
COPY --from=build /streamfair_identity_provider/identity_provider .

COPY sh ./sh
COPY db/migration ./db/migration

EXPOSE 8081
EXPOSE 9091

CMD [ "/streamfair_identity_provider/identity_provider" ]
ENTRYPOINT [ "/streamfair_identity_provider/start.sh" ]

RUN apk add --no-cache bash curl