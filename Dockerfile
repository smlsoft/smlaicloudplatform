FROM golang:1.17.3-alpine as builder

RUN apk add build-base
WORKDIR /go/src/app
ADD . /go/src/app

RUN go mod download
RUN go mod verify
RUN go build -a -o app .



FROM alpine
RUN apk add librdkafka
Add ./.env .env
COPY --from=builder /go/src/app/app /app
ENTRYPOINT ["./app"]