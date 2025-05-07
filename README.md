# Task Manager CLI

This repository provides a simple command-line task manager application written in Go. It manages tasks with basic operations like adding and listing them, and marking tasks as done. It’s organized into multiple packages and includes test coverage reporting in Cobertura format.

## Prerequisites

- Go 1.23+ (or a fairly recent version of Go)
- `gocover-cobertura` tool for Cobertura coverage conversion

## Installation

Make sure you have Go and `gocover-cobertura` installed:

```bash
go install github.com/t-yuki/gocover-cobertura@latest
```

## Building

To build the `taskmgr` binary:

```bash
go build ./cmd/taskmgr
```

This will create a `taskmgr` executable in the current directory.

## Usage

```bash
./taskmgr [command] [arguments]
```

### Commands

- `add <title>`: Adds a new task with the given title.
- `list`: Lists all tasks.
- `done <index>`: Marks the task at `<index>` as done.

Example:

```bash
./taskmgr add "Buy milk"
./taskmgr list
./taskmgr done 0
```

## Running Tests

To run all tests:

```bash
go test ./...
```

## Generating Coverage Report

To generate a Cobertura format coverage report:

```bash
go test ./... -coverprofile=coverage.out && gocover-cobertura < coverage.out > coverage.xml
```

This creates `coverage.xml` which you can use with CI tools that support Cobertura reports.
