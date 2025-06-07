FROM golang:alpine AS build
WORKDIR /auth
COPY auth-service/go.mod auth-service/go.sum ./
COPY auth-service ./
RUN go build -o app ./cmd

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /auth/app .
CMD [ "app" ]
