FROM golang:alpine AS build
WORKDIR /gateway
COPY gateway/go.mod gateway/go.sum ./
COPY gateway ./
RUN go build -o app ./cmd

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /gateway/app .
CMD [ "app" ]
