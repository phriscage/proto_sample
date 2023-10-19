# Database Files

* [Overview](#overview)
* [Data Model](#data-model)
* [Versioning](#versioning)

# Overview

This directory conatins all the relevant database information and files.


# Data Model

The data model is primarily mapped to the protobuf definitions and source of truth for data tables and fields.

## Protobuf sync

There are a few OSS libraries that will generate ORM mappings from protobuf, but my initial research has not yielded a stable library or tool to convert proto to SQL DDL statements for versioning. For now, this is a manual process outlined below.


# Versioning

Database versioning utilizes [Golang Migrate](https://github.com/golang-migrate) to handle migrations and versioning for table creation and modifications.

## Install

Install Migrate CLI via appropriate [Installation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation) option e.g. *Go Unversioned* with appropriate database tags

Prerequisites:

    sudo apt-get install gcc

Installation:

    CGO_ENABLED=1 go install -tags 'postgres' -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

## Configure

Create the migration files in the [migrations](./migrations/) directory with *sql* file extension based on the inital data model. e.g.

    migrate create -ext sql -dir migrations -seq persons
    migrate create -ext sql -dir migrations -seq books


## Migrate

Migrate the database migration files based on the database storage engine.

sqlite3

    mkdir -p data/sqlite

    migrate -source file://migrations/sqlite -database sqlite3://data/sqlite/data.db up

PostgreSQL


## Validate

Validate the schema is correct

sqlite3

    sqlite3 data/sqlite/data.db '.schema'

PostgreSQL


# Samples Queries

Run a sample query

sqlite3

    sqlite3 data/sqlite/data.db '.mode json' "SELECT * FROM books;

PostgreSQL
