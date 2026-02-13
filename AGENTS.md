# Agents Guide

This repository is a project for building a tmux session manager. Keep guidance simple and consistent with current conventions in the codebase.

## Quick Reference

## Build, Lint, Test

There is an app binary `bin/tm` that is built with `make`. Use `go test` to compile and run tests.

### Build / Compile

- `go test ./...`
  - Compiles all packages and runs tests.

### Lint / Static Checks

- `go vet ./...`
  - Default Go vet only. No extra linters configured.
- `gofmt -w .`
  - Use before commit to normalize formatting.

### Run All Tests

- `go test ./...`
  - Runs all unit tests.
- `go test ./... -count=1`
  - Bypass test caching when validating changes.

### Run a Single Test

- `go test ./... -run TestValid`
  - Run one test by name across packages.
- `go test ./ -run TestValid`
  - Run one test in the root package.
- `go test ./... -run TestRangeLength`
  - Example for other tests.

### Run Tests in One File

- `go test ./... -run TestRangeLength`
  - Go does not run by file directly; use `-run` with a test name.

## Notes

- Prefer clarity over cleverness.
- Keep public API stable unless explicitly requested.
