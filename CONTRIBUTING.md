# Contributing to gomb

Thanks for taking the time to contribute!

## Getting started

```bash
git clone https://github.com/ernlel/gomb.git
cd gomb
go test ./...
```

No external dependencies are needed for the core library (`go mod tidy` is a no-op).

## Structure

```
gomb.go              ← core library
gomb_test.go         ← core tests
pkg/html/            ← named element constructors (separate module)
pkg/markup_to_gomb/  ← HTML-to-Go converter (separate module)
cmd/gen-html/        ← code generator for pkg/html/elements_gen.go
examples/            ← runnable example programs
```

## Running tests

```bash
go test ./...                  # core
cd pkg/html && go test ./...   # html package
cd pkg/markup_to_gomb && go test ./...  # converter
```

## Code generation

The 114 named element constructors in `pkg/html/elements_gen.go` are generated:

```bash
go run ./cmd/gen-html
```

After changing the element list or naming rules in `cmd/gen-html/main.go`, regenerate and commit the result.

## Conventions

- `gomb` package: zero external dependencies. Core library only.
- Sub-packages in `pkg/` are separate Go modules with their own `go.mod`.
- Examples in `examples/` each have their own `go.mod` with `replace` directives pointing to the local checkout.
- Methods return a new `Element` (immutable value type). Never mutate in place.
- All text and attribute values are HTML-escaped unless passed through `Raw()`.
- Attribute keys sort alphabetically on render for deterministic output.

## Changelog

Keep entries in `CHANGELOG.md` under `[Unreleased]` until a release is cut.
When ready, change `[Unreleased]` to a version like `[1.0.0]`, and the CI
workflow will create the tag and release automatically on push.
