# Proto Sample

* [Overview](#overview)
* [Architecture](#architecture)
* [Quick Start](#quick-start)
* [Development](./DEVELOPMENT.md)
* [Data Model](#data-model)

# Overview

This Proto Sample repository provides a psuedo type-driven, data-model-first example for application design that showcases standardization and automatic generation of dependent data and application components. By starting with a well-defined, core data model using protobuf, application owners are able to create compiled libraries, sample servers or clients, and documentation using specific plugins and extensions based on the core model. This initial, non-exhaustive list includes the primary features:

* Protobuf data model
* gRPC service interface
* language specific libraries
* simple gRPC server
* CLI client examples

These additional features in the Proto Sample are utilized for simplcity:

* relational database storage
* object relational mapping

The [Data Model Dependency Graph](#data-model-dependency-graph) highlights data model dependencies with the foundational layer at the bottom which is utilized to automate the build-out of the upstream components.

## Data Model Dependency Graph

```mermaid
flowchart BT
  subgraph data-model[" "]
    direction LR
    DATA-MODEL([Data Model])
    INTERFACES([RPC Interfaces])
    PROTO([Proto Messages])
    DATA-MODEL -.->|uses|PROTO & INTERFACES
  end

  subgraph libraries[" "]
    direction RL
    LIBRARIES([Compiled Libraries])
    OPENAPI([OpenAPI Spec])
    DOCS([Documentation])
    STUBS([Client & Server Stubs])
    PROXY([Generated Proxy])
    SCHEMA([DB Schemas])
  end

  PROTO --> |protoc-gen|LIBRARIES
  PROTO --> |protoc-gen-sql/avro|SCHEMA
  INTERFACES --> |protoc-gen-grpc-gateway|PROXY
  INTERFACES --> |protoc-gen-openapi|OPENAPI
  INTERFACES --> |protoc-gen-docs|DOCS
  INTERFACES --> |protoc-gen-grpc|STUBS

  subgraph clients[" "]
    direction RL
    %% Golang & Python & JAVA & XYZ["..."]
    R-CLIENT["REST client"]
    G-CLIENT["gRPC client"]
  end

  subgraph docs[" "]
    G-DOCS["gRPC docs"]
    R-DOCS["REST docs"]
  end

  subgraph apis[" "]
    direction RL
    REST[REST Service] & GRPC[gRPC Service]
    REST <-->|transcodes|GRPC
  end

  PROXY -.->|imports|REST

  OPENAPI -.->|generates|R-DOCS
  OPENAPI -.->|utilizes|R-CLIENT

  DOCS -.-> |utilizes|G-DOCS
  LIBRARIES -.->|imports|G-CLIENT
  STUBS -.->|utilizes|G-CLIENT

  R-CLIENT -->|calls|REST
  G-CLIENT -->|calls|GRPC
  STUBS -..->|utilizes|GRPC

%% linkStyle 0,3 color:green;
class libraries,data-model,apis,clients,docs someclass;
classDef someclass fill:#F4FAFC,stroke:#333,stroke-width:2px,color:#fff,stroke-dasharray: 5 5
```

The gRPC service's interface, method, and REST transcoded resources follow the Google Cloud API and API Improvement Proposals naming standards ([here](https://cloud.google.com/apis/design/naming_convention) and [here](https://google.aip.dev/))


# Architecture


# Quick Start


## Validation

After you have an instance of the gRPC server running, you can use either the [grpcurl](#grpcurl-client) or the [samplectl](#samplectl-client) executable CLI client to validate and communicate with the gRPC server interface methods. Both clients leverage the same gRPC server interfeace API methods defined in the [sample_service.proto](./proto/sample/v1alpha/sample_service.proto)

### grpcurl client

**tl;dr** From [grpcurl](https://github.com/fullstorydev/grpcurl):
> grpcurl is a command-line tool that lets you interact with gRPC servers. It's basically curl for gRPC servers.

Install the `grpcurl` CLI to your local machine to communicate with the gRPC server. eg. OSX

    brew install grpcurl

List available gRPC server services (via reflection)

    grpcurl -plaintext localhost:10000 list

List all methods of an available service(s) (via reflection)

_**Note**_ You will need to import the service proto files (and import the directories of the proto dependencies) for the gRPC server service reflection to work. Download via [3P Proto Dependencies](#3p-proto-dependencies)

    grpcurl -plaintext -import-path third_party/googleapis -import-path third_party/protoc-gen-gorm -import-path proto -proto sample/v1alpha/sample_service.proto localhost:10000 list sample.v1alpha.SampleService

#### Sample interface methods

Try out some of the gRPC interface methods with `grpcurl`:

GetBook:

    grpcurl -plaintext -d '{"name": "123"}' -import-path third_party/googleapis -import-path third_party/protoc-gen-gorm -import-path proto -proto sample/v1alpha/sample_service.proto localhost:10000 sample.v1alpha.SampleService/GetBook

CreateBook:

    grpcurl -plaintext -d '{"book": {"name": "123"} }' -import-path third_party/googleapis -import-path third_party/protoc-gen-gorm -import-path proto -proto sample/v1alpha/sample_service.proto localhost:10000 sample.v1alpha.SampleService/CreateBook

ListBooks:

    grpcurl -plaintext -d '{"name_prefix": "1234"}' -import-path third_party/googleapis -import-path third_party/protoc-gen-gorm -import-path proto -proto sample/v1alpha/sample_service.proto localhost:10000 sample.v1alpha.SampleService/ListBooks


## Database

Follow instructions in the [db/README](db/README.md)

### 3P proto dependencies

Download protobuf 3P import dependencies for `grpcurl`

    mkdir -p third_party

    git clone https://github.com/googleapis/googleapis third_party/googleapis
    git clone https://github.com/infobloxopen/protoc-gen-gorm third_party/protoc-gen-gorm


# Wishlist

These are some additional feature components I would like to include/showcase in this example

* ~~automated build tool~~ (Bazel added)
* deployment configurations
* gRPC to REST transcoding
* custom data model field extensions
* schema and DDL converters
* tests, tests, tests!
