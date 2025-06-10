FROM golang:alpine AS build
WORKDIR /notification
COPY notification-service/go.mod notification-service/go.sum ./
COPY notification-service ./
RUN go build -o app ./cmd

FROM alpine:3.21
WORKDIR /usr/bin
COPY --from=build /notification/app .
EXPOSE 8080
CMD ["app"]