# Proto Sample

* [Overview](#overview)
* [Architecture](#architecture)
* [Quick Start](#quick-start)
* [Development](#development)
* [Data Model](#data-model)

# Overview

This Proto Sample repository showcases a sample protobuf data model and generating language libraries and CLI binary executable client.


# Architecture


# Quick Start

## Setup

## Client


# Development

Install the Proto Sample executable on a local development machine to test the functionality and generate the appropriate libraries.

First, verify you have golang >= 1.20.x installed, or download from [Go.dev](https://go.dev/dl/)

    $ go version
    go version go1.20 linux/amd64

Next, clone or download this project and download the package dependencies

    go mod download

If initial version, instatiate `go mod` and `go mod tidy`

    go mod init github.com/phriscage/proto_sample
    go mod tidy


## Protobuf

Generate the protobuf files for the following languages: *Golang*

Install Protobuf tools and libraries:
*OSX*

    brew install protobuf

*Linux*

    sudo apt-get install -y protobuf-compiler
    #sudo apt-get install -y golang-goprotobuf-dev


Install Go plugins:

    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    go install github.com/infobloxopen/protoc-gen-gorm@latest

Download protobuf 3P build dependencies: 
_TODO_ Add [Buf](https://buf.build/) to streamline proto build process
*Make directories*

    mkdir -p third_party

*Service*

    git clone https://github.com/googleapis/googleapis third_party/googleapis

*GORM* 

    git clone https://github.com/infobloxopen/protoc-gen-gorm third_party/protoc-gen-gorm

*BQ Schemas*

    git clone https://github.com/GoogleCloudPlatform/protoc-gen-bq-schema third_party/protoc-gen-bq-schmea

Generate libraries:
*Make directories*

    mkdir -p gen/go

*Service*

    protoc -I proto -I third_party/googleapis -I third_party/protoc-gen-gorm --go_out ./gen/go/ --go_opt paths=source_relative --go-grpc_out ./gen/go/ --go-grpc_opt paths=source_relative proto/sample/v1alpha/*.proto

*GROM* 

    protoc -I proto -I third_party/googleapis -I third_party/protoc-gen-gorm --gorm_out ./gen/go/ --gorm_opt paths=source_relative proto/sample/v1alpha/*.proto

*BQ Schemas*

    protoc -I temp -I third_party/protoc-gen-bq-schema --bq-schema_out=temp/bq_schema temp/bq_schema/foo.proto


## Clients

Generate the CLI binary executable client:

Create the [/bin](./bin) directory and export in the PATH environment variable:

    mkdir -p bin && export PATH=$PATH:${PWD}/bin


Build the client

    GO111MODULE=on go build -o bin/samplectl ./cmd/sample

Test

    samplectl -h


