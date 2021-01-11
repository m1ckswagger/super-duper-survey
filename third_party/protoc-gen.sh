#!/bin/bash
protoc --proto_path=api/proto/v1 --proto_path=third_party \
    --go_out=pkg/api/v1 --go_opt=paths=source_relative \
    --go-grpc_out=pkg/api/v1 --go-grpc_opt=paths=source_relative \
    question.proto
    # catalog.proto
