# Testing

* [Prequisistes](#prequisistes)
* [gRPCurl](#grpcurl)
    * [Examples](#grpcurl-examples)
* [Proto Sample Control](#proto-sample-control)
    * [Examples](#proto-sample-control-examples)

This document includes details for how to send sample messages and payloads to the [Proto Sample](./README.md) server interfaces. You can use either the [gRPCurl](#grpcurl) or the [Proto Sample Control](#proto-sample-control) command line interface (CLI) utility to validate and communicate with the gRPC server interface methods. Both clients leverage the same gRPC server interfeace API methods defined in the [sample_service.proto](./proto/sample/v1alpha/sample_service.proto)

_**Note**_ For these examples, I am using an *OSX* development environment (arm64, Apple M2) with [Homebrew](https://brew.sh/) installed.

## Prequisistes

A compiled and running instance of the Proto Sample gRPC server (e.g [Development#executables](./DEVELOPMENT.md#executables) that is reachable from your test environment.


## grpcurl

**tl;dr** From [grpcurl](https://github.com/fullstorydev/grpcurl):
> grpcurl is a command-line tool that lets you interact with gRPC servers. It's basically curl for gRPC servers.

_**Note**_ You will need to import the service proto files (and import the directories of the proto dependencies) for the gRPC server service reflection to work. Download via [3P proto dependencies](#3p-proto-dependencies)


Install the `grpcurl` CLI to your local machine to communicate with the gRPC server. eg. OSX

    brew install grpcurl

List available gRPC server services (via reflection)

    grpcurl -plaintext localhost:10000 list

List all methods of an available service(s) (via reflection)

    grpcurl -plaintext -import-path third_party/googleapis -import-path third_party/protoc-gen-gorm -import-path proto -proto sample/v1alpha/sample_service.proto localhost:10000 list sample.v1alpha.SampleService

### grpcurl client examples

Try out some of the gRPC interface methods examples with `grpcurl`:

GetBook:

    grpcurl -plaintext -d '{"name": "123"}' -import-path third_party/googleapis -import-path third_party/protoc-gen-gorm -import-path proto -proto sample/v1alpha/sample_service.proto localhost:10000 sample.v1alpha.SampleService/GetBook

CreateBook:

    grpcurl -plaintext -d '{"book": {"name": "123"} }' -import-path third_party/googleapis -import-path third_party/protoc-gen-gorm -import-path proto -proto sample/v1alpha/sample_service.proto localhost:10000 sample.v1alpha.SampleService/CreateBook

ListBooks:

    grpcurl -plaintext -d '{"name_prefix": "1234"}' -import-path third_party/googleapis -import-path third_party/protoc-gen-gorm -import-path proto -proto sample/v1alpha/sample_service.proto localhost:10000 sample.v1alpha.SampleService/ListBooks

## Proto Sample Control

The [Proto Sample Control](./cmd/sample) (a.k.a. `samplectl`) CLI is a Golang utility that can be utilized to execute commands against the Proto Sample gRPC server and its interfaces.

_**Note**_ You will need to compile an instance of the Proto Sample Control CLI Golang binary via [Development#executables](./DEVELOPMENT.md#executables).


### Proto Sample Control examples

TBD


## 3P proto dependencies

Download protobuf 3P import dependencies for `grpcurl`

    mkdir -p third_party

    git clone https://github.com/googleapis/googleapis third_party/googleapis
    git clone https://github.com/infobloxopen/protoc-gen-gorm third_party/protoc-gen-gorm

