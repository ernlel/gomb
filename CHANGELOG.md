# Changelog

All notable changes to gomb will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] — 2026-07-08

### Added
- `.Class(names ...string)` — shorthand for `.A("class", Classes(...))`.
- `.Id(id string)` — shorthand for `.A("id", id)`.
- `.When(cond bool, fn func(*Element))` — chainable conditional transform method. Unlike package-level `When()`, this stays inside the builder chain.
- `.Clone()` — shallow copy with independent Attributes map. ChildNodes slice references are shared.
- `ErrNilWriter` — exported sentinel error returned by `Render(nil)` instead of `io.ErrShortWrite`.

### Changed
- `A(pairs ...string)` — now accepts variadic key-value pairs. Backward-compatible: existing `.A(k, v)` calls work unchanged. Odd trailing arguments are silently ignored.
- `Render(w)` — now returns `(int64, error)` instead of just `error`. Useful for logging bytes written or checking write errors.

### Fixed
- `Render(nil)` returns `ErrNilWriter` (identifiable sentinel) instead of `io.ErrShortWrite`.

## [1.1.0] — 2026-07-08

### Changed
- **Breaking:** Switched from value receivers to pointer receivers (`*Element`). `E()`, `A()`, `T()`, `C()` and all helpers now return `*Element`. Methods mutate in place and return the same pointer for chaining — no more accidental copy discards.
- `If(cond, el)`, `IfElse(cond, a, b)`, `When(cond, fn)`, `Map(slice, fn)`, `Fragment(els...)` all return `*Element`.
- `None` is now `nil` (`var None *Element`). `ToString()` and `Render()` handle nil receivers gracefully.
- `With(fns...)` transformers now take `func(*Element)` (mutate in place).
- `C()` skips nil children.
- Immutability copy-on-write removed — `A()` and `C()` mutate the element directly.

## [1.0.0] — 2026-07-07

### Added
- `Classes(...names)` — space-join class names, skipping empties, for safe conditional CSS.
- `.Data(key, value)` — shorthand for `data-*` attributes.
- `.Style(css)` — shorthand for the `style` attribute.
- `.Attrs(pairs...)` / `.As(pairs...)` — apply multiple `Attr` pairs at once.
- `NS(prefix)` — create a namespaced attribute builder (e.g. `hx := NS("hx-")`).
- `.With(fns...)` — apply composable transformer functions to an element.
- `When(cond, fn)` — lazy conditional: `fn` is only called if `cond` is true.
- `El(tag)`, `.Attr(k,v)`, `.Text(s)`, `.Children(...)` — long-form aliases for `E/A/T/C`.
- `Txt(s)` — tag-less text element constructor.
- Exported struct fields `Attributes`, `ChildNodes`, `TextContent` for introspection.
- `pkg/html` sub-package — 114 named HTML element constructors with inline-style API.
- `pkg/markup_to_gomb` — converts HTML markup strings to gomb Go code.
- `cmd/gen-html` — code generator for `pkg/html/elements_gen.go`.
- `.github/workflows/test.yml` — CI: tests, lint, and example builds.
- `.golangci.yml` — linter configuration.
- `LICENSE` — MIT.

### Changed
- `Attr` struct field renamed to `Attributes`.
- `Children` field renamed to `ChildNodes`.
- `Text` field renamed to `TextContent`.
- Attribute keys are sorted alphabetically on render (deterministic output).
- Self-closing tag lookup changed from O(n) slice scan to O(1) map lookup.
- Project layout: `html/` → `pkg/html/`, generators → `cmd/`, converter → `pkg/markup_to_gomb/`.
- Root `go.mod` has zero external dependencies (`x/net` moved to `pkg/markup_to_gomb`).

### Fixed
- `Render()` no longer panics on nil `io.Writer` — returns an error.
- `Paragraph` example component: fragile `Children = append` replaced with `Map` + `Fragment`.

### Security
- `.T()`, `.A()`, `Txt()` now HTML-escape all values — XSS-safe by default.
- `<script>` and `<style>` text content is never entity-encoded.
- `Raw()` provides explicit opt-in for unescaped content.

[1.2.0]: https://github.com/ernlel/gomb/releases/tag/v1.2.0
[1.1.0]: https://github.com/ernlel/gomb/releases/tag/v1.1.0
[1.0.0]: https://github.com/ernlel/gomb/releases/tag/v1.0.0
