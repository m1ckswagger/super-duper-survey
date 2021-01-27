# Survey

**All third party code used in this package is located in subfolder `third_party`**

## Setting up the go
1. Init the project

    ```bash
    go mod init <package_name>
    ```

2. Create folders for API definition

    ```bash
    mkdir -p api/proto/v1
    ```

3. Create the file `catalog.proto`

4. Download `protoc` 

    Downloading binaries from the [official repo](https://github.com/protocolbuffers/protobuf/releases) (rather do not use package manager as versions are most likely dated).

5. Get the **Go plugins** for `protoc`

    ```bash
    $ go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
    ```

6. Ensure to update `PATH` for `protoc` to find plugins

    ```bash
    $ export PATH="$PATH:$(go env GOPATH)/bin"
    ```

7. Add proto build script `./third_party/protoc-gen.sh`

    ```bash
    #!/bin/bash
    protoc --proto_path=api/proto/v1 --proto_path=third_party --go_out=plugins=grpc:pkg/api/v1 todo-service.proto
    ```