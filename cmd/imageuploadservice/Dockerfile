
## smlsoft/smlaicloudplatform:imagesservice
##
## Builder
##
FROM golang:1.17-alpine3.15 AS builder

RUN apk add alpine-sdk
RUN apk add librdkafka=1.8.2-r0
RUN apk add build-base

WORKDIR /go/app

COPY go.mod /go/app
COPY go.sum /go/app
RUN go mod download

ADD . /go/app

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build  -o go-app -tags musl ./cmd/imageuploadservice/main.go  

##
## Deploy
##
FROM alpine:latest 
WORKDIR /root/
# ADD ./cmd/authenticationservice/.env .
ENV SERVICE_NAME images-api

COPY --from=builder /go/app/go-app .
ENTRYPOINT  /root/go-app