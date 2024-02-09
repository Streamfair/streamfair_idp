# Build Stage
FROM golang:1.22.0-alpine3.19 AS build
WORKDIR /streamfair_idp
COPY . .
RUN go build -o streamfair_idp main.go

# Run Stage
FROM alpine:3.19
WORKDIR /streamfair_idp
COPY --from=build /streamfair_idp/streamfair_idp .
COPY app.env .

EXPOSE 8081
EXPOSE 9091

CMD ["./streamfair_idp"]