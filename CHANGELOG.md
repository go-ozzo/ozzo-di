# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

No version has been tagged yet; the first tagged release will be v1.0.0.

## [Unreleased]

### Added
- `go.mod` with the module path `github.com/go-ozzo/ozzo-di` (the import path is unchanged) and Go 1.21 as the minimum. The package has no dependencies
- CI on GitHub Actions: one required check `Test` that gathers every gated job; tests with `-race` on Go oldstable and stable on Linux, stable on Windows and macOS, and the `go.mod` minimum; a linux/386 build and test; pinned golangci-lint v2.15.0 with a strict `.golangci.yml`; `gofmt` and `go mod tidy` checks; govulncheck (informational); coverage to Codecov; Dependabot for GitHub Actions
- `SECURITY.md`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `.gitattributes`

### Changed
- No change in library behaviour. Mechanical edits only:
  - a test printed an `int` with `%q`; it now uses `%d` (message only). `go vet` flags this in modules that declare Go 1.26 or later
  - doc comments reformatted by `gofmt` (Go 1.19 style)
  - `reflect.Pointer` and `reflect.PointerTo` instead of the older `reflect.Ptr` and `reflect.PtrTo` names, a type switch with assignment, and a preallocated result slice in `Call` (still a non-nil slice)
- The README requirements say Go 1.21 or above; the Travis CI config and its badges are removed
