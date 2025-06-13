FROM golang:alpine AS build
WORKDIR /settlement
COPY settlement/go.mod settlement/go.sum ./ 
COPY settlement ./
RUN go build -o app ./cmd

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /settlement/app .
CMD [ "app" ]