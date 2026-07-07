# Changelog

All notable changes to gomb will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
- `Attr` struct field renamed to `Attributes` (exported, non-conflicting with `.Attrs()` method).
- `Children` field renamed to `ChildNodes` (exported, for introspection).
- `Text` field renamed to `TextContent` (exported, for introspection).
- `.A()` now deep-copies the attribute map (true immutability).
- `.C()` now deep-copies the children slice (true immutability).
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
