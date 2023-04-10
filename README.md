# Golang with Gin

A showcase written in go utilizing the Gin web framework for delivery and fx framework for dependency injection.

Project specification: [SPEC.md](../../main/SPEC.md)

This project is **not** production ready. The following concerns should be addressed to make this project production ready:

- Internal models are presented to the end user. \
  Usecase methods should return data through models prepared for presentation
- Internal models serve as input models. \
  Controllers should utilize models prepared for data input instead of internal models
- Audit log might fail. \
  End-user might receive a non-success response because the audit log failed to create an entry.
  This behaviour might cause confusion or undefined behaviour on the client side. \
  There are several ways to deal with this issue:
  - Insert audit in a transaction alongside the main database interaction
  - Handle audit asynchronously
  - Repeat audit action until success
- No database migrations. \
  There is nothing in place to migrate model changes to database tables and table creation is statically written
- No proper SQL builder. \
  Although any SQL query syntax errors will be raised before the server starts, queries can be written incorrectly. \
  Implementation of any more advanced SQL features might require a lot of changes
- Errors are not verbose enough. \
  Error messages returned to the end user aren't verbose enough which might make debugging harder
- Only the most basic tests are written. \
  Tests only test the most basic and successful situation.
  They should be more exhaustive
- This project was written with multiple databases in mind. \
  Although some production-ready servers might support multiple databases, generally this approach adds a lot of unnecessary requirements and constraints

## Build

This project build like any golang project

```console
$ go build
```

This command will generate the `clean-show` executable file

## Test

Before you test, you need to install [vektra/mockery](https://github.com/vektra/mockery) and generate mock files

Run these commands to test this project:

```console
$ go install github.com/vektra/mockery/v2@v2.20.0
$ go generate ./...
$ go test ./...
```
