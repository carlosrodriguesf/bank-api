FROM golang:1.18-alpine

WORKDIR /opt/app/api

RUN apk update
RUN apk add git openssh

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

ENV GO_PRIVATE="github.com/carlosrodriguesf"

CMD [ "go", "run", "pkg/main.go" ]