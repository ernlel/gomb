# gomb — Go Markup Builder

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-00ADD8)](https://go.dev)
[![Test](https://github.com/ernlel/gomb/actions/workflows/test.yml/badge.svg)](https://github.com/ernlel/gomb/actions/workflows/test.yml)
[![License](https://img.shields.io/badge/License-MIT-green)](./LICENSE)

`gomb` is a small Go library for building HTML programmatically using a fluent,
type-safe API. There are no templates, no string concatenation, and no reflection —
just regular Go functions and method chaining.

```
go get github.com/ernlel/gomb
```

## Why gomb?

| Pain point with text templates | How gomb solves it |
|---|---|
| Syntax errors only found at runtime | Go compiler checks everything |
| Hard to compose and reuse fragments | Components are plain Go functions |
| Escaping easy to forget | `T()` and `A()` HTML-escape automatically |
| Attribute order is non-deterministic in templates | gomb sorts attributes alphabetically |
| No IDE refactoring support | Full Go tooling: rename, find-references, … |

## Quick start

```go
import . "github.com/ernlel/gomb"   // dot-import makes E, If, Map, etc. available unqualified

page := E("html").A("lang", "en").C(
    E("head").C(
        E("meta").A("charset", "UTF-8"),
        E("title").T("My App"),
    ),
    E("body").C(
        E("h1").T("Hello, World!"),
        E("p").T("Welcome to gomb."),
    ),
)

fmt.Println(page.ToString())
```

Output:

```html
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>
      My App
    </title>
  </head>
  <body>
    <h1>
      Hello, World!
    </h1>
    <p>
      Welcome to gomb.
    </p>
  </body>
</html>
```

Every short method has a long-form alias. Use whichever reads better in context:

```go
// Short (terse DSL)
page := E("html").A("lang", "en").C(
    E("head").C(E("title").T("My App")),
    E("body").C(E("h1").T("Hello"), E("p").T("Welcome")),
)

// Long (self-documenting)
page := El("html").Attr("lang", "en").Children(
    El("head").Children(El("title").Text("My App")),
    El("body").Children(El("h1").Text("Hello"), El("p").Text("Welcome")),
)

// Named inline (closest to HTML)
page := Html(
    Head(TitleElement(Txt("My App"))),
    Body(H1(Txt("Hello")), P(Txt("Welcome"))),
).Attr("lang", "en")
```

| Short | Long |
|-------|------|
| `E(tag)` | `El(tag)` |
| `.A(pairs...)` | `.Attr(pairs...)` |
| `.T(text)` | `.Text(text)` |
| `.C(elems...)` | `.Children(elems...)` |

### Named element constructors

For maximum IDE autocomplete and compile-time safety, the `html` sub-package
provides named constructor functions for every HTML element. Import alongside
`gomb`:

```go
import (
    . "github.com/ernlel/gomb"      // core API
    . "github.com/ernlel/gomb/pkg/html"  // named constructors
)
```

Use the **inline style** for compact, HTML-like code — pass attributes and
children directly as arguments:

```go
page := Html(
    Head(
        Meta(Attr{Key: "charset", Value: "UTF-8"}),
        TitleElement(Txt("My App")),
    ),
    Body(
        H1(Txt("Hello, World!")),
        P(Txt("Welcome to gomb.")),
    ),
).Attr("lang", "en")
```

Or use the **chained style** with explicit `.Attr()` / `.Children()` / `.Text()`:

```go
page := Html().Attr("lang", "en").Children(
    Head().Children(
        Meta().Attr("charset", "UTF-8"),
        TitleElement().Text("My App"),
    ),
    Body().Children(
        H1().Text("Hello"),
        P().Text("Welcome to gomb."),
    ),
)
```

Both produce identical output. `Txt(s)` is shorthand for `E("").T(s)` —
a text node without a wrapping tag.

The `html` package is a separate module (`github.com/ernlel/gomb/pkg/html`) and
depends only on `gomb`. Import it only when you want named constructors; the
core package stays minimal.

All 114 standard elements — `Anchor()` to `Wbr()`. A handful use an `Element`
suffix to avoid colliding with existing gomb functions:

| HTML tag | Constructor | Reason |
|---|---|---|
| `<a>` | `Anchor()` | Collides with `.A()` / `.Attr()` setters |
| `<input>` | `InputElement()` | Common variable name |
| `<script>` | `ScriptElement()` | Common variable name |
| `<style>` | `StyleElement()` | Collides with `.Style()` helper |
| `<title>` | `TitleElement()` | Common variable name |
| `<data>` | `DataElement()` | Collides with `.Data()` helper |
| `<map>` | `MapElement()` | Collides with `Map[T]()` |
| `<var>` | `VarElement()` | Collides with `var` keyword |
| `<time>` | `TimeElement()` | Collides with `time` package |

Regenerate with: `go run ./cmd/gen-html`

## Core API

| Function / method | Description |
|---|---|
| `E(tag)` / `El(tag)` | Create an element |
| `.A(pairs...)` / `.Attr(pairs...)` | Set attributes from key-value pairs (HTML-escaped) |
| `.T(text)` / `.Text(text)` | Set text content (HTML-escaped) |
| `.C(children...)` / `.Children(children...)` | Append child elements |
| `.Class(names...)` | Shorthand for the `class` attribute (uses `Classes()`) |
| `.Id(id)` | Shorthand for the `id` attribute |
| `.Clone()` | Shallow copy — independent Attributes map |
| `.When(cond, fn)` | Apply `fn` to element only when `cond` is true |
| `.ToString()` | Render to an HTML string |
| `.Render(w)` | Write HTML to an `io.Writer`, returns `(int64, error)` |
| `Raw(html)` | Insert pre-rendered HTML verbatim |
| `Fragment(els...)` | Tag-less wrapper (no extra element) |
| `None` | Empty element — renders nothing |
| `If(cond, el)` | Conditionally include an element |
| `IfElse(cond, a, b)` | Choose between two elements |
| `When(cond, fn)` | Lazy conditional — `fn` only called if `cond` is true |
| `Map(slice, fn)` | Transform a slice into `[]*Element` |
| `Classes(...names)` | Space-join class names, skipping empties |
| `.Data(key, value)` | Shorthand for `data-*` attributes |
| `.Style(css)` | Shorthand for the `style` attribute |
| `.Attrs(pairs...)` / `.As(pairs...)` | Apply multiple `Attr` pairs at once |
| `NS(prefix)` | Create a namespaced attribute builder |
| `.With(fns...)` | Apply composable transformer functions |
| `Div()`, `Span()`, … | Named constructors in `github.com/ernlel/gomb/pkg/html` |
| `Txt(text)` | Tag-less text node — shorthand for `E("").T(text)` |

## HTML escaping

`T()` and `A()` escape `&`, `<`, `>`, and `"` automatically — user-supplied
data is safe by default.

```go
E("p").T(`<script>alert("xss")</script>`)
// → <p>
//     &lt;script&gt;alert(&#34;xss&#34;)&lt;/script&gt;
//   </p>
```

`<script>` and `<style>` content is **never** entity-encoded (browsers parse
these tags as raw text).

```go
E("script").T(`if (a < b) console.log("ok")`)
// → <script>
//     if (a < b) console.log("ok")
//   </script>
```

Use `Raw()` to inject pre-rendered fragments or external HTML:

```go
E("div").C(Raw("<span>already safe HTML</span>"))
```

## Fragments

`Fragment(els...)` creates a **tag-less wrapper** — it groups multiple elements
together without emitting any surrounding HTML tag. This is useful when a
component needs to return several sibling elements, or when you want to insert
multiple items into a parent without an extra `<div>`.

```go
// Without Fragment you'd need a wrapper div:
E("div").C(
    E("dt").T("Go"),
    E("dd").T("A statically typed language."),
)

// Fragment emits the children directly — no wrapper tag:
func definition(term, desc string) *gomb.Element {
    return Fragment(
        E("dt").T(term),
        E("dd").T(desc),
    )
}

dl := E("dl").C(
    definition("Go",   "A statically typed language."),
    definition("gomb", "An HTML builder for Go."),
)
```

Output:

```html
<dl>
  <dt>
    Go
  </dt>
  <dd>
    A statically typed language.
  </dd>
  <dt>
    gomb
  </dt>
  <dd>
    An HTML builder for Go.
  </dd>
</dl>
```

Another common use: returning a group of table cells or list items from a
helper without wrapping them in an extra element:

```go
func navLinks(links []Link) *gomb.Element {
    items := Map(links, func(l Link) *gomb.Element {
        return E("li").C(E("a").A("href", l.URL).T(l.Label))
    })
    return Fragment(items...)   // no wrapping <ul>/<div>
}

E("ul").C(navLinks(siteLinks))
```

## Named attribute spaces

`NS(prefix)` creates a function that prepends the given prefix to every attribute key.
Combine with `.Attrs()` to keep htmx, Alpine.js, and custom `data-*` attributes tidy
without repeating the prefix on every call.

### htmx

```go
hx := gomb.NS("hx-")
E("button").Attrs(
    hx("get", "/api/users"),
    hx("swap", "outerHTML"),
    hx("target", "#user-list"),
    hx("trigger", "click"),
).T("Reload")
```

Output:

```html
<button hx-get="/api/users" hx-swap="outerHTML" hx-target="#user-list" hx-trigger="click">
  Reload
</button>
```

### Alpine.js

```go
x := gomb.NS("x-")
E("div").Attrs(
    x("data", `{"open": false}`),
    x("show", "open"),
    x("cloak", ""),
)
```

### Custom data attributes

```go
data := gomb.NS("data-")
E("li").Attrs(
    data("user", "42"),
    data("role", "admin"),
    data("sort", "name"),
)
// → <li data-user="42" data-role="admin" data-sort="name">
```

### Mixing namespaces with plain attributes

`Attrs()` works alongside `.A()` and `.Data()`:

```go
hx := gomb.NS("hx-")
E("input").
    A("type", "text").
    A("class", "form-input").
    Data("validate", "email").
    Attrs(
        hx("get", "/validate"),
        hx("trigger", "keyup changed delay:500ms"),
    )
```

## CSS helpers

### .Class()

`.Class(names...)` is a shorthand for `.A("class", Classes(names...))`:

```go
E("button").Class(
    "btn",
    IfElse(isPrimary, "btn-primary", "btn-secondary"),
    IfElse(isLarge, "text-lg", "text-sm"),
)
// → <button class="btn btn-primary text-lg">
```

### Classes()

`Classes()` joins class names into a single string, skipping empty strings. Useful when
you need the string value outside of `.Class()`:

```go
E("button").A("class", Classes(
    "btn",
    IfElse(isPrimary, "btn-primary", "btn-secondary"),
    IfElse(isDisabled, "opacity-50 cursor-not-allowed", ""),
))
```

### .Id()

`.Id(id)` is a shorthand for `.A("id", id)`:

```go
E("div").Id("main").Class("container")
// → <div class="container" id="main">
```

### Style()

```go
E("div").Style("display: flex; align-items: center; gap: 0.5rem")
// → <div style="display: flex; align-items: center; gap: 0.5rem">
```

## Element transformers

`.With()` accepts one or more `func(*Element)` transformers and applies them in order.
Package common attribute patterns as reusable, composable modifiers:

```go
// Transformer factories — functions that return transformers
btn := func(e *Element) {
    e.A("type", "submit").A("class", Classes("btn", "rounded"))
}
primaryBtn := func(e *Element) {
    e.A("class", Classes("btn", "btn-primary", "rounded"))
}
withTooltip := func(text string) func(*Element) {
    return func(e *Element) {
        e.Data("tooltip", text)
    }
}

// Compose them
E("button").T("Save").With(btn, withTooltip("Save your changes"))
E("button").T("Delete").With(primaryBtn, withTooltip("Permanently delete"))

// Transformers can also modify children:
withIcon := func(e *Element) {
    e.C(E("span").A("class", "icon").T("★"))
}
E("a").A("href", "/star").T("Favorite").With(withIcon)
```

Transformers are just functions — they're testable, composable, and can be defined
in packages alongside your components.

## Conditional chaining

`.When(cond, fn)` applies `fn` to the element only when `cond` is true. Unlike the
package-level `When()` (which lazily builds elements for `.C()`), this method stays
inside the chain:

```go
E("input").
    A("type", "checkbox").
    When(checked, func(e *Element) { e.A("checked", "") }).
    When(disabled, func(e *Element) { e.A("disabled", "") })
```

## Cloning elements

`.Clone()` returns a shallow copy — the Attributes map is independent, but ChildNodes
slice references are shared. Use it to build a base template and stamp out variants:

```go
baseInput := E("input").Class("input").A("type", "text")

nameInput := baseInput.Clone().A("name", "name").Id("name-field")
emailInput := baseInput.Clone().A("name", "email").Id("email-field")
// baseInput is unchanged
```

## Rendering

`Render(w)` writes directly to any `io.Writer` and returns `(int64, error)` —
useful for logging bytes written or checking write errors:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    n, err := E("p").T("Hello").Render(w)
    log.Printf("wrote %d bytes", n)
}
```

## Components

Components are Go functions that return `*Element`. Any parameters, loops, and
conditions are just Go:

```go
func Button(label string, primary bool) *Element {
    return E("button").
        A("class", Classes("btn", IfElse(primary, "btn-primary", ""))).
        A("type", "submit").
        T(label)
}

func NavBar(links []NavLink) *Element {
    items := Map(links, func(l NavLink) *Element {
        return E("li").C(E("a").A("href", l.URL).T(l.Label))
    })
    return E("nav").C(E("ul").C(items...))
}
```

## Layouts

Wrap content in a layout function:

```go
func Layout(title string, body *gomb.Element) *gomb.Element {
    return E("html").A("lang", "en").C(
        E("head").C(
            E("meta").A("charset", "UTF-8"),
            E("title").T(title),
            E("link").A("rel", "stylesheet").A("href", "/static/app.css"),
        ),
        E("body").C(Header(), body, Footer()),
    )
}

func IndexPage() *gomb.Element {
    return Layout("Home", E("main").C(
        E("h1").T("Welcome"),
        E("p").T("This is the home page."),
    ))
}
```

## Conditionals and loops

```go
// If — render an element only when a condition is true
E("ul").C(
    If(user.IsAdmin, E("li").C(E("a").A("href", "/admin").T("Admin"))),
    E("li").C(E("a").A("href", "/profile").T("Profile")),
)

// IfElse — choose between two elements
IfElse(user.LoggedIn,
    E("a").A("href", "/logout").T("Log out"),
    E("a").A("href", "/login").T("Log in"),
)

// Go if/else inside a component function
func greeting(user *User) *gomb.Element {
    if user == nil {
        return E("p").T("Hello, guest!")
    }
    return E("p").T("Hello, " + user.Name + "!")
}

// Map — render a list from a slice
names := []string{"Alice", "Bob", "Carol"}
E("ul").C(Map(names, func(name string) *gomb.Element {
    return E("li").T(name)
})...)

// When — lazy conditional: fn is only called if the condition is true
// Use When instead of If when the element is expensive or has side-effects
E("ul").C(
    When(user != nil, func() *Element {
        // This closure only runs when user != nil
        return E("li").C(E("a").A("href", "/profile").T(user.Name))
    }),
    E("li").C(E("a").A("href", "/home").T("Home")),
)
```

## Tailwind CSS

Tailwind classes are regular strings — use `Classes()` for terse conditional
composition:

```go
E("button").A("class", Classes(
    "bg-blue-600 hover:bg-blue-700",
    "text-white font-bold",
    "py-2 px-4 rounded",
    IfElse(disabled, "opacity-50 cursor-not-allowed", ""),
)).T("Click me")
```

For Tailwind CDN (development / demos):

```go
E("link").
    A("rel", "stylesheet").
    A("href", "https://cdn.jsdelivr.net/npm/tailwindcss@3/dist/tailwind.min.css")
```

## htmx integration

htmx directives are HTML attributes — use `.A()` directly, or define a namespace
with `NS()` for less repetition:

```go
// Vanilla A() calls
E("button").
    A("hx-get", "/api/users").
    A("hx-target", "#user-list").
    A("hx-swap", "innerHTML").
    T("Reload users")

// Or with NS() — no repeating "hx-"
hx := NS("hx-")
E("button").Attrs(
    hx("get", "/api/users"),
    hx("target", "#user-list"),
    hx("swap", "innerHTML"),
).T("Reload users")

// An inline-edit form using namespace
E("form").Attrs(
    hx("put", fmt.Sprintf("/users/%d", user.ID)),
    hx("target", fmt.Sprintf("#user-%d", user.ID)),
    hx("swap", "outerHTML"),
).C(
    E("input").A("type", "text").A("name", "name").A("value", user.Name),
    E("button").A("type", "submit").T("Save"),
)
```

See [`examples/htmx/`](examples/htmx/main.go) for a full task-list app.

## Alpine.js integration

Alpine.js `x-data`, `x-show`, `@click`, `:class`, and other directives are just
HTML attributes. JSON in `x-data` is HTML-escaped automatically and the browser
decodes it back correctly.

```go
x := NS("x-")

// Counter
E("div").Attrs(x("data", `{"count": 0}`)).C(
    E("button").A("@click", "count--").T("-"),
    E("span").A("x-text", "count"),
    E("button").A("@click", "count++").T("+"),
)

// Toggle visibility
E("div").Attrs(
    x("data", `{"open": false}`),
    x("show", "open"),
    x("cloak", ""),
).C(
    E("button").A("@click", "open = !open").T("Toggle"),
    E("p").A("x-show", "open").T("Visible when open is true"),
)
```

See [`examples/alpinejs/`](examples/alpinejs/main.go) for accordion, tabs, modal,
and live search examples.

## Caching

Rendered HTML strings can be cached and replayed freely.

**Static component (sync.Once)**

```go
var (
    navOnce sync.Once
    navHTML string
)

func Nav() string {
    navOnce.Do(func() { navHTML = buildNav().ToString() })
    return navHTML
}
```

**Per-key TTL cache**

```go
var cache sync.Map

type entry struct { html string; exp time.Time }

func Cached(key string, ttl time.Duration, build func() *gomb.Element) string {
    if v, ok := cache.Load(key); ok {
        if e := v.(entry); time.Now().Before(e.exp) {
            return e.html
        }
    }
    html := build().ToString()
    cache.Store(key, entry{html, time.Now().Add(ttl)})
    return html
}
```

See [`examples/caching/`](examples/caching/main.go) for a complete demo including
an HTTP response-level cache middleware.

## Converting existing HTML

The `pkg/markup_to_gomb` package converts an HTML string into the equivalent gomb
Go code. Useful for migrating existing templates or pasting in HTML from a
design tool.

```go
import "github.com/ernlel/gomb/pkg/markup_to_gomb"

code, err := markup_to_gomb.GenerateGombFromMarkup(`
    <div class="card">
        <h2>Title</h2>
        <p>Body text.</p>
    </div>
`)
fmt.Println(code)
```

See [`examples/markup_to_gomb/`](examples/markup_to_gomb/main.go) for a
file-based conversion tool.

## Examples

| Example | What it shows |
|---|---|
| [`examples/components/`](examples/components/) | Reusable component functions served over HTTP |
| [`examples/layout/`](examples/layout/) | Multi-page site with a shared layout, header, and footer (Tailwind) |
| [`examples/htmx/`](examples/htmx/) | Dynamic task list with htmx partial updates |
| [`examples/alpinejs/`](examples/alpinejs/) | Client-side interactivity (counter, accordion, tabs, modal, live search) |
| [`examples/caching/`](examples/caching/) | Component cache, per-key TTL cache, and response-level cache middleware |
| [`examples/test1/`](examples/test1/) | Full landing page generated from gomb (Tailwind) |
| [`examples/html-elements/`](examples/html-elements/) | Named constructors for all 114 HTML elements — inline and chained styles |
| [`examples/markup_to_gomb/`](examples/markup_to_gomb/) | Convert HTML markup to gomb code |

## Self-closing tags

The following void elements are rendered without a closing tag automatically:
`area`, `base`, `br`, `col`, `embed`, `hr`, `img`, `input`, `link`, `meta`,
`param`, `source`, `track`, `wbr`.

```go
E("img").A("src", "logo.png").A("alt", "Logo")
// → <img alt="Logo" src="logo.png" />
```

## Boolean attributes

Pass an empty string value to render a boolean attribute without a value:

```go
E("input").A("type", "checkbox").A("checked", "").A("disabled", "")
// → <input checked disabled type="checkbox" />
```

## Introspection

The `Element` struct fields are exported so you can inspect or walk the tree
in custom tooling, tests, or middleware:

```go
el := E("div").A("class", "card").C(
    E("h2").T("Title"),
    E("p").T("Body"),
)

// Inspect
tag := el.Tag                 // "div"
cls := el.Attributes["class"] // "card"
txt := el.ChildNodes[0].TextContent  // "Title"

// Walk the tree
var walk func(e *Element)
walk = func(e *Element) {
    fmt.Println(e.Tag)
    for _, c := range e.ChildNodes {
        walk(c)
    }
}
walk(el)
```

Builder methods (`A`, `C`, `T`, etc.) mutate the element in place and return
the same pointer for chaining — read any field at any time.

## Code style guide

### Short or long API — choose by audience

The short form (`E`, `A`, `T`, `C`) reads like a DSL and is ideal for page-level
composition where you're writing lots of markup:

```go
page := E("html").A("lang", "en").C(
    E("head").C(E("title").T("Home")),
    E("body").C(E("h1").T("Welcome")),
)
```

The long form (`El`, `Attr`, `Text`, `Children`) is self-documenting — use it in
public packages, shared components, and code that non-Go developers might read:

```go
func SiteLayout(title string, body *Element) *Element {
    return El("html").Attr("lang", "en").Children(
        El("head").Children(El("title").Text(title)),
        El("body").Children(Header(), body, Footer()),
    )
}
```

### Chain vertically for clarity

One method call per line — attributes grouped, children indented:

```go
// Good
E("button").
    A("type", "submit").
    A("class", Classes("btn", "btn-primary")).
    Data("action", "save").
    T("Save")

// Hard to scan
E("button").A("type","submit").A("class",Classes("btn","btn-primary")).Data("action","save").T("Save")
```

### Components are functions

Extract anything used more than once into a function. Return `*Element`, accept any
parameters you need:

```go
func Card(title, body string) *Element {
    return E("div").A("class", "card").C(
        E("h3").A("class", "card-title").T(title),
        E("p").A("class", "card-body").T(body),
    )
}

func CardList(items []CardData) *Element {
    return E("div").A("class", "card-list").C(
        Map(items, func(d CardData) *Element {
            return Card(d.Title, d.Body)
        })...,
    )
}
```

### Reuse attribute bundles with NS()

When a prefix repeats across many elements, define a namespace at the package or
function level:

```go
var hx = gomb.NS("hx-")

func liveSearch() *Element {
    return E("input").As(
        hx("get", "/search"),
        hx("trigger", "keyup changed delay:300ms"),
        hx("target", "#results"),
        hx("swap", "innerHTML"),
    ).A("type", "search").A("placeholder", "Search...")
}
```

### Parameterize transformers

`func(*Element)` is the signature — wrap in a closure to accept
configuration:

```go
func sizedText(size string) func(*Element) {
    return func(e *Element) {
        e.A("class", Classes("text-"+size))
    }
}

E("p").T("Title").With(sizedText("xl"))
E("p").T("Body").With(sizedText("base"))
```

### Conditionals — pick the right tool

| Pattern | When |
|---|---|
| Go `if/else` inside a component | The condition drives logic, not just presence |
| `If(cond, el)` | Inline inside a `.C()` call, element already built |
| `When(cond, fn)` | The element is expensive or has side-effects |
| `IfElse(cond, a, b)` | Swapping between two inline elements |
| `Classes("base", IfElse(cond, "on", "off"))` | Conditional CSS classes |

```go
// Go if/else — multiple branches or complex logic
func StatusBadge(status string) *Element {
    var color string
    switch status {
    case "active":   color = "bg-green-500"
    case "paused":   color = "bg-yellow-500"
    default:         color = "bg-gray-500"
    }
    return E("span").A("class", Classes("badge", color)).T(status)
}

// If — simple inline presence
E("ul").C(
    If(isAdmin, E("li").C(E("a").A("href", "/admin").T("Admin"))),
    E("li").C(E("a").A("href", "/home").T("Home")),
)

// Classes + IfElse — conditional CSS
E("button").A("class", Classes(
    "btn",
    IfElse(active, "btn-active", "btn-inactive"),
    IfElse(large, "text-lg", ""),
))
```

### Layouts: wrap, don't repeat

A layout function accepts `title` and `body`, returns the full page:

```go
func Layout(title string, body *Element) *Element {
    return E("html").A("lang", "en").C(
        E("head").C(
            E("meta").A("charset", "UTF-8"),
            E("title").T(title),
        ),
        E("body").C(Header(), body, Footer()),
    )
}
```

Every page becomes one line:

```go
func HomePage() *Element     { return Layout("Home", homeContent()) }
func AboutPage() *Element    { return Layout("About", aboutContent()) }
```

## License

MIT — see [LICENSE](./LICENSE). Contributions welcome — see [CONTRIBUTING.md](./CONTRIBUTING.md).

