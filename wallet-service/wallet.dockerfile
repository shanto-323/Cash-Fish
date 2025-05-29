FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /wallet-service
COPY go.mod go.sum ./
COPY wallet-service wallet-service
RUN go build -o app ./wallet-service/cmd

FROM alpine:3.21
WORKDIR /usr/bin
COPY --from=build /wallet-service/app .
EXPOSE 8080
CMD ["app"]
