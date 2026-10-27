# Contributing to ozzo-di

Thank you for your interest in contributing to ozzo-di! This document covers how to build, test, and submit changes.

## Prerequisites

- **Go 1.21+** ([download](https://go.dev/dl/))

## Building and Testing

```bash
go build ./...
go test -race ./...
```

`-race` needs cgo and a C compiler; without one, run `go test ./...`. CI runs the race detector on Linux, Windows
and macOS.

## Code Style

- Run `gofmt -w .` and `golangci-lint run --timeout=5m` before every commit. CI enforces both.
- Run `go vet ./...` and fix all issues.
- Follow standard Go naming conventions (`ID`, `URL`, `HTTP` are uppercase).
- Exported types and functions must have doc comments.
- The package has no third-party dependencies; please keep it that way.

## Compatibility

Code imports this package today, so the exported API and its behaviour change only on purpose: a behaviour change
needs a test that shows it and an entry under `Changed` in `CHANGELOG.md`. A bug fix starts with a test that fails
without it.

## Pull Request Workflow

1. Fork the repository and create a feature branch:
   ```bash
   git checkout -b fix/my-fix
   ```
2. Make your changes. Keep commits focused and well-described.
3. Verify locally:
   ```bash
   go build ./...
   go vet ./...
   gofmt -l .
   golangci-lint run --timeout=5m
   go test -race ./...
   ```
4. Add a line to the `[Unreleased]` section of `CHANGELOG.md` for any user-visible change.
5. Push and open a pull request against `master`.
6. Wait for CI to pass. All checks must be green before merge.

Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/):
```
fix: resolve dependencies of a parent binding in the child container
docs: fix the inject tag in the README example
test: cover providers registered as shared
```

## Reporting Bugs

Open an [issue](https://github.com/go-ozzo/ozzo-di/issues) with a clear description, a minimal program that shows
the problem, and the Go version you're using. Security issues go through the Security tab instead (see
[SECURITY.md](SECURITY.md)).

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
