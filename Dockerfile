FROM golang

ADD . /go/src/survey

WORKDIR /go/src/survey/cmd/server

RUN go mod tidy
RUN go build

ENTRYPOINT ./server

EXPOSE 433
