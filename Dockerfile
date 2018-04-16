FROM golang:latest
RUN mkdir -p /go/src/github.com/whiteshtef/p
ADD .  /go/src/github.com/whiteshtef/p
WORKDIR  /go/src/github.com/whiteshtef/p

RUN go get -v
RUN go build -o main .
EXPOSE 80

CMD ["/go/src/github.com/whiteshtef/p/main","--logtostderr"]
