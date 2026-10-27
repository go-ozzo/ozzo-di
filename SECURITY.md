# Security Policy

## Supported Versions

No version has been tagged yet. Until the first release (v1.0.0), fixes land on `master`.
After that, security fixes are made for the latest release.

## Reporting a Vulnerability

**DO NOT** open a public GitHub issue for security vulnerabilities.

Report them privately through the repository's Security tab:
https://github.com/go-ozzo/ozzo-di/security/advisories/new

### What to Include

- Description of the vulnerability
- Steps to reproduce
- Affected versions (commit or pseudo-version)
- Potential impact

### Response Timeline

- **Initial Response**: Within 72 hours
- **Fix & Disclosure**: Coordinated with reporter

## Security Considerations

ozzo-di builds values and calls functions through reflection. Users should be aware of:

1. **Types come from code, not from input.** The container builds whatever type it is asked for and calls every
   provider and every function passed to `Call()`. Never let untrusted input choose the types, providers or
   functions the container uses.
2. **Dependency cycles.** A struct whose `inject`-tagged fields lead back to its own type (directly or through other
   types) makes `Make()`, `Inject()` and `Call()` recurse until the goroutine's stack overflows, which ends the
   process. Keep the dependency graph acyclic.
3. **Concurrency.** A container has no locking. Register everything before the container is shared between
   goroutines. Resolving a type bound to a shared provider stores the provider's result in the container on first
   use, so concurrent `Make()`, `Inject()` or `Call()` on such a type is a data race as well.
4. **Panics.** `RegisterAs()` panics when the value or type cannot be converted to the target type, and
   `InterfaceOf()` panics when given a non-interface, and `Call()` when given a non-function or `nil`. A provider
   must return a valid `reflect.Value` of the registered type: an invalid one makes `Make()`, `Inject()` and
   `Call()` panic, and one of another type makes `Inject()` and `Call()` panic, while `Make()` returns it
   unchecked. A parent container must be one created by `NewContainer()`: resolving through any other
   `Container` implementation set with `SetParentContainer()` panics. Recover where a configuration error must not stop the program.

## Security Contact

- **GitHub Security Advisory**: https://github.com/go-ozzo/ozzo-di/security/advisories/new
