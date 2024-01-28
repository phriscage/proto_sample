# Development

* [Golang](#golang)
* [Libraries](#libraries)
    * [Buf](#buf)
    * [Protoc](#protoc)
* [Executables](#executables)
    * [Bazel](#bazel)
    * [Go](#go)
    * [Docker](#docker)
* [Database](#database)

Setup a local development environment to generate the Proto Sample [libraries][#libraries] & [executables][#executables] to test and develop the functionality of the appropriate services. There are automated (preferred) and manual instrauctions to build and compile each component and appropriate dependencies. Installing [Golang][#golang] is a prerequisite for either installation method.

_**Note**_ For these examples, I am using an *OSX* development environment (arm64, Apple M2) with [Homebrew](https://brew.sh/) installed.

## Golang

First, verify you have golang >= 1.20.x installed, or download from [Go.dev](https://go.dev/dl/)

    go version

Set the GO_PATH environment variables in you shell profile config after install:

    echo "export GO_PATH=~/go" >> ~/.bash_profile
    echo "export PATH=$PATH:/$GO_PATH/bin" >> ~/.bash_profile
    source ~/.bash_profile

Next, clone or download this project and download the package dependencies

    go mod download

If initial version, instatiate `go mod` and `go mod tidy`

    go mod init github.com/phriscage/proto_sample
    go mod tidy

## Libraries

The protobuf [files](./proto) contain the source for the client & server libraries and generated stubs. See [data-model-dependency-graph](./README.md#data-model-dependency-graph). [Protoc](https://grpc.io/docs/protoc-installation/) is the protobuf compiler, [protoc-gen-X](https://protobuf.dev/reference/go/go-generated/) are the plugin extensions, and [Buf](https://buf.build) is used for managing protobuf packages and dependencies. Buf is the automated way for building and Protoc can be executed manually.

Install Protobuf compliler, Protoc, and libraries:

    brew install protobuf

Install the protoc-gen-go, protoc-gen-go-grpc, protoc-gen-gorm protoc plugin(s) for the output language of choice (Go, Go gRPC, GORM respectively) and set your $PATH. These plugins are used for either Buf or Protoc build options:

    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
    go install github.com/infobloxopen/protoc-gen-gorm@latest

Export PATH env with GOPATH:

    export PATH="$PATH:$(go env GOPATH)/bin"

### Buf

Install the [Buf CLI](https://buf.build/docs/installation) for managing protobuf dependencies and generating libraries.

    brew install bufbuild/buf/buf

Download and cache protobuf dependencies:

    buf mod update proto/

_**Note**_: The proto-gen-X plugins and binaries are defined in the [buf.gen.yaml](./buf.gen.yaml) configuration

Generate the libraries and client/server stubs:

    buf generate

### Protoc

Generate libraries: Make directories

    mkdir -p gen/go

Service

    protoc -I proto -I third_party/googleapis -I third_party/protoc-gen-gorm --go_out ./gen/go/ --go_opt paths=source_relative --go-grpc_out ./gen/go/ --go-grpc_opt paths=source_relative proto/sample/v1alpha/*.proto

GROM

    protoc -I proto -I third_party/googleapis -I third_party/protoc-gen-gorm --gorm_out ./gen/go/ --gorm_opt paths=source_relative proto/sample/v1alpha/*.proto


## Executables

Generate the gRPC server and executable CLI client.

### Bazel

[Bazel](https://bazel.build/) is an open-source build and test tool similar to Make, Maven, and Gradle based on Google's internal Blaze tool

_**Note**_: For now, Bazel utilizes the pre-build proto [libraries](#libraries) above vs generating them with it's plugins

Install the Bazel package via Homebrew as follows:

    brew install bazel

Run the Bazel Gazelle plugin to generate the BUILD files in the appropriate directories.

    bazel run //:gazelle

Build the server

    bazel build //server

Build the client

    bazel build //cmd/sample

Test both

    bazel-bin/server/server_/server -h

    bazel-bin/cmd/sample/sample_/sample -h

### Go

Create the [/bin](./bin) directory and export in the PATH environment variable:

    mkdir -p bin && export PATH=$PATH:${PWD}/bin


Build the server

    GO111MODULE=on go build -o bin/samplesrv ./server/.

Build the client

    GO111MODULE=on go build -o bin/samplectl ./cmd/sample

Test both

    samplesrv -h

    samplectl -h


### Docker

Docker container is another option to test during local development. These instructions showcase building the Docker container locally from [Dockerfile](./Dockerfile) using [Docker build](https://docs.docker.com/build/guide/) and running in Docker engine on your local machine.

Install [Docker Desktop or Docker Engine](https://docs.docker.com/get-docker/) and verify

    docker version

Build the container image

    docker build -t proto-sample-server:dev .

Run the container with the local file system and DB and DEBUG log severity

    docker run --rm -it -v ${PWD}/db:/app/db -p 10000:10000 proto-sample-server:dev -host 0.0.0.0 -log_severity DEBUG -database_connection_dsn db/data/sqlite/data.db

Test with either `grpcurl` or executable client above

    grpcurl -plaintext localhost:10000 list


## Database

Follow the instructions in the [db/README](db/README.md) to setup a database instance.

