#!/bin/bash
protoc --proto_path=api/proto/v1 --proto_path=third_party \
    --go_out=pkg/api/v1 --go_opt=paths=source_relative \
    --go-grpc_out=pkg/api/v1 --go-grpc_opt=paths=source_relative \
    catalog.proto

protoc --proto_path=api/proto/v1 --proto_path=third_party \
    --go_out=pkg/api/v1 --go_opt=paths=source_relative \
    --go-grpc_out=pkg/api/v1 --go-grpc_opt=paths=source_relative \
    question.proto

protoc --proto_path=api/proto/v1 --proto_path=third_party \
    --grpc-gateway_out=logtostderr=true:pkg/api/v1 catalog.proto

protoc --proto_path=api/proto/v1 --proto_path=third_party \
    --swagger_out=logtostderr=true:api/swagger/v1 catalog.proto