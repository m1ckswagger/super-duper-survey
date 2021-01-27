#!/bin/bash

protos="catalog user answer"

for proto in $protos; do
    protoc --proto_path=api/proto/v1 --proto_path=third_party \
        --go_out=pkg/api/v1 --go_opt=paths=source_relative \
        --go-grpc_out=pkg/api/v1 --go-grpc_opt=paths=source_relative \
        $proto.proto

    protoc --proto_path=api/proto/v1 --proto_path=third_party \
        --grpc-gateway_opt logtostderr=true \
        --grpc-gateway_opt paths=source_relative \
        --grpc-gateway_out=logtostderr=true:pkg/api/v1 $proto.proto

    protoc --proto_path=api/proto/v1 --proto_path=third_party \
        --swagger_out=logtostderr=true:api/swagger/v1 $proto.proto
done