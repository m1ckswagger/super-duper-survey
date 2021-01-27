FROM golang

ADD . /go/src/survey

WORKDIR /go/src/survey

RUN go mod tidy

# ENTRYPOINT go run main.go -grpc-port=9090 -http-port=8080 -db-host=127.0.0.1 -db-user=root -db-password=survey -db-schema=survey

WORKDIR /go/src/survey/cmd/server

EXPOSE 8080
EXPOSE 9090
