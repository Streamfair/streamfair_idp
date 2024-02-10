# Build Stage
FROM golang:1.22.0-alpine3.19 AS build
WORKDIR /streamfair_identity_provider
COPY . .
RUN go build -o streamfair_identity_provider main.go

# Run Stage
FROM alpine:3.19
WORKDIR /streamfair_identity_provider
COPY --from=build /streamfair_identity_provider/streamfair_identity_provider .
COPY app.env .

EXPOSE 8081
EXPOSE 9091

CMD ["./streamfair_identity_provider"]