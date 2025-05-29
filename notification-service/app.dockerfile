FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /notification-service
COPY go.mod go.sum ./
COPY notification-service notification-service
RUN go build -o app ./notification-service/cmd

FROM alpine:3.21
WORKDIR /usr/bin
COPY --from=build /notification-service/app .
EXPOSE 8080
CMD ["app"]