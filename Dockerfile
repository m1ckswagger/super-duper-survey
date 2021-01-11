FROM golang

ADD . /go/src/survey

WORKDIR /go/src/survey

RUN go mod tidy
RUN go build

ENTRYPOINT ./survey

EXPOSE 433
