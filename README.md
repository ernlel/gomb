# gomb — Go Markup Builder

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

## Core API

| Function / method | Description |
|---|---|
| `E(tag)` | Create an element |
| `.A(key, value)` | Set an attribute (HTML-escaped) |
| `.T(text)` | Set text content (HTML-escaped) |
| `.C(children...)` | Append child elements |
| `.ToString()` | Render to an HTML string |
| `.Render(w)` | Write HTML to an `io.Writer` |
| `Raw(html)` | Insert pre-rendered HTML verbatim |
| `Fragment(els...)` | Tag-less wrapper (no extra element) |
| `None` | Empty element — renders nothing |
| `If(cond, el)` | Conditionally include an element |
| `IfElse(cond, a, b)` | Choose between two elements |
| `Map(slice, fn)` | Transform a slice into `[]Element` |

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
func definition(term, desc string) gomb.Element {
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
func navLinks(links []Link) gomb.Element {
    items := Map(links, func(l Link) gomb.Element {
        return E("li").C(E("a").A("href", l.URL).T(l.Label))
    })
    return Fragment(items...)   // no wrapping <ul>/<div>
}

E("ul").C(navLinks(siteLinks))
```



`Render(w)` writes directly to any `io.Writer`, including `http.ResponseWriter`:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    E("p").T("Hello").Render(w)
}
```

## Components

Components are Go functions that return `Element`. Any parameters, loops, and
conditions are just Go:

```go
func Button(label, href string, primary bool) gomb.Element {
    cls := "btn"
    if primary {
        cls += " btn-primary"
    }
    return E("a").A("href", href).A("class", cls).T(label)
}

func NavBar(links []NavLink) gomb.Element {
    items := Map(links, func(l NavLink) gomb.Element {
        return E("li").C(E("a").A("href", l.URL).T(l.Label))
    })
    return E("nav").C(E("ul").C(items...))
}
```

## Layouts

Wrap content in a layout function:

```go
func Layout(title string, body gomb.Element) gomb.Element {
    return E("html").A("lang", "en").C(
        E("head").C(
            E("meta").A("charset", "UTF-8"),
            E("title").T(title),
            E("link").A("rel", "stylesheet").A("href", "/static/app.css"),
        ),
        E("body").C(Header(), body, Footer()),
    )
}

func IndexPage() gomb.Element {
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
func greeting(user *User) gomb.Element {
    if user == nil {
        return E("p").T("Hello, guest!")
    }
    return E("p").T("Hello, " + user.Name + "!")
}

// Map — render a list from a slice
names := []string{"Alice", "Bob", "Carol"}
E("ul").C(Map(names, func(name string) gomb.Element {
    return E("li").T(name)
})...)
```

## Tailwind CSS

Tailwind classes are regular attribute values — just pass them to `.A("class", ...)`:

```go
E("button").
    A("class", "bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded").
    T("Click me")
```

For Tailwind CDN (development / demos):

```go
E("link").
    A("rel", "stylesheet").
    A("href", "https://cdn.jsdelivr.net/npm/tailwindcss@3/dist/tailwind.min.css")
```

## htmx integration

htmx directives are HTML attributes — use `.A()` exactly as you would for any
other attribute. The server returns HTML fragments for partial page updates.

```go
// A button that loads /api/users and swaps the content of #user-list
E("button").
    A("hx-get", "/api/users").
    A("hx-target", "#user-list").
    A("hx-swap", "innerHTML").
    T("Reload users")

// An inline-edit form
E("form").
    A("hx-put", fmt.Sprintf("/users/%d", user.ID)).
    A("hx-target", fmt.Sprintf("#user-%d", user.ID)).
    A("hx-swap", "outerHTML").
    C(
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
// Counter
E("div").
    A("x-data", `{"count": 0}`).
    C(
        E("button").A("@click", "count--").T("-"),
        E("span").A("x-text", "count"),
        E("button").A("@click", "count++").T("+"),
    )

// Toggle visibility
E("div").
    A("x-data", `{"open": false}`).
    C(
        E("button").A("@click", "open = !open").T("Toggle"),
        E("p").A("x-show", "open").T("Visible when open is true"),
    )
```

See [`examples/alpinejs/`](examples/alpinejs/main.go) for accordion, tabs, modal,
and live search examples.

## Caching

Because `Element` is an immutable value type, rendered HTML strings can be stored
and replayed freely.

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

func Cached(key string, ttl time.Duration, build func() gomb.Element) string {
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

The `tools/markup_to_gomb` package converts an HTML string into the equivalent gomb
Go code. Useful for migrating existing templates or pasting in HTML from a
design tool.

```go
import "github.com/ernlel/gomb/tools/markup_to_gomb"

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

## License

MIT

