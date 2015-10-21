FROM golang:1.5.1
MAINTAINER Brian Jacobel <brian@bjacobel.com>

ADD . /go/src/github.com/bjacobel/slact
RUN go install github.com/bjacobel/slact

ENTRYPOINT /go/bin/slact

EXPOSE 3000
