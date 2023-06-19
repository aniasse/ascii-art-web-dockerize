FROM golang:alpine

LABEL MAINTAINER="Fatidiop Aboubakdiallo Aniasse" VERSION="1.0"

WORKDIR /app

ADD . /app

RUN go build -o main .

RUN apk update && apk add bash && apk add tree

EXPOSE 8080

CMD ["/app/main"]
