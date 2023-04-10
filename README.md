# Golang with Gin

A showcase written in go utilizing the Gin web framework for delivery and fx framework for dependency injection.

Project specification: [SPEC.md]()

This project is **not** production ready. The following concorns should be addressed to make this project production ready:

- Internal models are presented to the end user. \
  Usecase methods should return data through models prepared for presentation
- Internal models serve as input models. \
  Controller should utilize models prepared for data input instead of internal models
- Audit log might fail resulting in a failed request. \
  End-user might receive non-success response because audit log failed to create an entry.
  This behavior might couse confusion or undifined behaviour on the client side. \
  There are several ways to deal with this issue:
  - Insert audit in a transaction alogside the main database interaction
  - Handle audit asynchronously
  - Repeat audit action unitl success
- No database migrations. \
  There is nothing in place to migrate model changes to database tables and table creation is statically written
- No proper SQL builder. \
  Although any SQL query errors will be rised before server start, queries can be written incorrectly
- Errors are not verbose eniugh. \
  Error messages returned to the end ser aren't verbose enough which might make debugging harder
- Only most basic test are written. \
  Tests only test the most basic and successful situation.
  They should be more exhaustive

## Build

This project build like any golang project

```console
$ go build
```

This command will generate the `clean-show` executable file

## Test

Before you test, you need to generate mock files

Run these commands to test this project:

```console
$ go generate ./...
$ go test ./...
```
