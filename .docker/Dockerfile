FROM golang:1.18-alpine

RUN apk update
RUN apk add git openssh make gcc libc-dev

RUN go install github.com/golang/mock/mockgen@v1.6.0

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

ENV GO_PRIVATE="github.com/carlosrodriguesf"

WORKDIR /opt/app/api

CMD [ "make", "go-run" ]