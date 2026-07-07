// Package gomb provides a fluent, type-safe API for building HTML programmatically in Go.
// Elements are immutable value types: every method returns a new Element, making them
// safe to share and reuse across goroutines without synchronisation.
//
// Quick start:
//
//	import . "github.com/ernlel/gomb"
//
//	page := E("html").A("lang", "en").C(
//	    E("head").C(E("title").T("My App")),
//	    E("body").C(
//	        E("h1").T("Hello, World!"),
//	        E("p").T("Welcome to gomb."),
//	    ),
//	)
//	fmt.Println(page.ToString())
package gomb

import (
	"html"
	"io"
	"sort"
	"strings"
)

// Element represents an HTML element. It is an immutable value type; all mutating
// methods return a new Element rather than modifying the receiver.
//
// Fields are exported for introspection and custom tooling (e.g. walking the tree,
// inspecting state). For building elements prefer the fluent API: E, A/T/C, etc.
type Element struct {
	Tag        string
	Attributes map[string]string
	ChildNodes []Element
	TextContent string
	rawText    bool // when true, TextContent is written verbatim (no HTML escaping)
}

// None is the zero Element. It renders to nothing and is the canonical return value
// when a conditional helper has no output.
var None = Element{}

// E creates a new HTML element with the given tag name.
// Use an empty tag ("") to create a tag-less fragment.
func E(tag string) Element {
	return Element{Tag: tag}
}

// El is an alias for E, for users who prefer a slightly longer name.
// Identical to E in every way.
func El(tag string) Element {
	return Element{Tag: tag}
}

// Txt creates a tag-less text element. It is equivalent to E("").T(text) —
// a text node without a wrapping tag. Useful with inline constructors from
// the html package:
//
//	Div(gomb.Attr{Key: "class", Value: "card"}, H2(Txt("Title")))
func Txt(content string) Element {
	return Element{TextContent: content}
}

// Raw creates an element whose content is written verbatim to the output without any
// HTML escaping. Use it for pre-rendered HTML fragments or for content inside
// <script> and <style> tags that must not be entity-encoded.
//
//	E("script").A("type", "text/javascript").C(Raw(`console.log("hello")`))
//	E("style").C(Raw(`.btn { color: red; }`))
//	E("div").C(Raw("<span>pre-rendered</span>"))
func Raw(rawHTML string) Element {
	return Element{TextContent: rawHTML, rawText: true}
}

// Fragment creates a tag-less wrapper around multiple elements. The fragment itself
// produces no markup; only its children are rendered.
//
//	E("ul").C(Fragment(
//	    E("li").T("one"),
//	    E("li").T("two"),
//	))
func Fragment(elements ...Element) Element {
	return E("").C(elements...)
}

// A returns a new Element with the given attribute set (or overwritten).
// The attribute value is HTML-escaped when the element is rendered, so it is safe
// to pass raw strings including characters like <, >, &, and ".
//
// Boolean attributes (e.g. "disabled", "checked") can be set by passing an empty value:
//
//	E("input").A("type", "checkbox").A("checked", "")
func (e Element) A(key, value string) Element {
	// Copy the attribute map so this element is independent of the original.
	newAttrs := make(map[string]string, len(e.Attributes)+1)
	for k, v := range e.Attributes {
		newAttrs[k] = v
	}
	newAttrs[key] = value
	e.Attributes = newAttrs
	return e
}

// Attr is an alias for A — identical behaviour, longer name for readability.
func (e Element) Attr(key, value string) Element {
	return e.A(key, value)
}

// T returns a new Element with the given text content. The text is HTML-escaped
// when the element is rendered, preventing XSS when displaying user-supplied data.
// For unescaped output use Raw().
func (e Element) T(value string) Element {
	e.TextContent = value
	e.rawText = false
	return e
}

// Text is an alias for T — identical behaviour, longer name for readability.
func (e Element) Text(value string) Element {
	return e.T(value)
}

// C returns a new Element with the given elements appended to its children.
// Passing None or a zero Element has no visible effect on the rendered output.
//
// To append a slice of elements use the spread operator:
//
//	items := Map(list, func(s string) Element { return E("li").T(s) })
//	E("ul").C(items...)
func (e Element) C(elements ...Element) Element {
	// Copy to avoid sharing the backing array with the original slice.
	newChildren := make([]Element, len(e.ChildNodes)+len(elements))
	copy(newChildren, e.ChildNodes)
	copy(newChildren[len(e.ChildNodes):], elements)
	e.ChildNodes = newChildren
	return e
}

// Children is an alias for C — identical behaviour, longer name for readability.
func (e Element) Children(elements ...Element) Element {
	return e.C(elements...)
}

// ToString renders the element tree to an indented HTML string.
func (e Element) ToString() string {
	var sb strings.Builder
	e.render(&sb, "")
	return sb.String()
}

// Render writes the HTML representation of the element to w.
// It is a convenience wrapper around ToString for use with http.ResponseWriter
// and other io.Writer targets.
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    w.Header().Set("Content-Type", "text/html; charset=utf-8")
//	    page.Render(w)
//	}
func (e Element) Render(w io.Writer) error {
	if w == nil {
		return io.ErrShortWrite
	}
	_, err := io.WriteString(w, e.ToString())
	return err
}

// Classes returns a space-joined string from the given class names, skipping
// empty strings. This makes conditional CSS classes ergonomic:
//
//	E("div").A("class", Classes("base", IfElse(active, "active", "")))
func Classes(names ...string) string {
	var sb strings.Builder
	for i, n := range names {
		if n != "" {
			if i > 0 && sb.Len() > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(n)
		}
	}
	return sb.String()
}

// Data returns a "data-*" attribute builder for data attributes.
//
//	E("div").Data("count", "0").Data("user", "42")
//	// <div data-count="0" data-user="42">
func (e Element) Data(key, value string) Element {
	return e.A("data-"+key, value)
}

// Style sets the "style" attribute on the element. For multiple CSS declarations
// pass a semicolon-separated string:
//
//	E("div").Style("color: red; font-weight: bold")
func (e Element) Style(css string) Element {
	return e.A("style", css)
}

// Attr is a key-value attribute pair.
type Attr struct {
	Key, Value string
}

// Attrs applies multiple Attr pairs to the element in one call.
//
//	E("div").Attrs(
//	    Attr{"class", "container"},
//	    Attr{"id", "main"},
//	)
func (e Element) Attrs(attrs ...Attr) Element {
	result := e
	for _, a := range attrs {
		result = result.A(a.Key, a.Value)
	}
	return result
}

// As is an alias for Attrs — the short plural form for applying multiple
// attribute pairs at once.
//
//	hx := gomb.NS("hx-")
//	E("button").As(hx("get", "/api"), hx("swap", "outerHTML")).T("Load")
func (e Element) As(attrs ...Attr) Element {
	return e.Attrs(attrs...)
}

// NS returns a function that builds namespaced attribute pairs. The returned
// function prepends the given prefix to every key, enabling domain-specific
// attribute shortcuts without repeating the prefix.
//
//	// htmx attributes
//	hx := gomb.NS("hx-")
//	E("button").Attrs(
//	    hx("swap", "outerHTML"),
//	    hx("target", "#result"),
//	    hx("get", "/api/data"),
//	)
//
//	// Alpine.js x-data / x-show / x-text etc.
//	x := gomb.NS("x-")
//	E("div").Attrs(x("data", `{"open": false}`), x("show", "open"))
//
//	// Custom data attributes
//	data := gomb.NS("data-")
//	E("li").Attrs(data("user", "42"), data("role", "admin"))
func NS(prefix string) func(key, value string) Attr {
	return func(key, value string) Attr {
		return Attr{Key: prefix + key, Value: value}
	}
}

// With applies the given transformer functions to the element and returns the
// result. Transformer functions receive and return an Element, allowing custom
// logic to be packaged as reusable modifiers.
//
//	// A transformer that adds common attributes
//	btn := func(e gomb.Element) gomb.Element {
//	    return e.A("type", "submit").A("class", "btn")
//	}
//	// A transformer that conditionally adds children
//	withHelp := func(e gomb.Element) gomb.Element {
//	    return e.C(E("span").T("?"))
//	}
//	E("form").C(
//	    E("button").T("Save").With(btn),
//	    E("input").With(btn, withHelp),
//	)
func (e Element) With(fns ...func(Element) Element) Element {
	result := e
	for _, fn := range fns {
		result = fn(result)
	}
	return result
}

// When returns el derived from fn() when cond is true, otherwise returns None.
// Unlike If, the element is not constructed unless the condition is true.
//
//	E("ul").C(
//	    When(user != nil, func() Element {
//	        return E("li").T("Hello, " + user.Name)
//	    }),
//	)
func When(cond bool, fn func() Element) Element {
	if cond {
		return fn()
	}
	return None
}

// If returns el when cond is true, otherwise returns None.
//
//	E("ul").C(
//	    If(isLoggedIn, E("li").C(E("a").A("href", "/logout").T("Logout"))),
//	    E("li").C(E("a").A("href", "/home").T("Home")),
//	)
func If(cond bool, el Element) Element {
	if cond {
		return el
	}
	return None
}

// IfElse returns ifEl when cond is true, otherwise returns elseEl.
//
//	IfElse(isLoggedIn,
//	    E("a").A("href", "/logout").T("Logout"),
//	    E("a").A("href", "/login").T("Login"),
//	)
func IfElse(cond bool, ifEl, elseEl Element) Element {
	if cond {
		return ifEl
	}
	return elseEl
}

// Map transforms a slice of values into a slice of Elements using fn.
// Use the spread operator to pass the result to C():
//
//	names := []string{"Alice", "Bob", "Carol"}
//	E("ul").C(Map(names, func(name string) Element {
//	    return E("li").T(name)
//	})...)
func Map[T any](items []T, fn func(T) Element) []Element {
	result := make([]Element, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}

// selfClosingTags is a set of void HTML elements that must not have a closing tag.
var selfClosingTags = map[string]bool{
	"area": true, "base": true, "br": true, "col": true,
	"embed": true, "hr": true, "img": true, "input": true,
	"link": true, "meta": true, "param": true, "source": true,
	"track": true, "wbr": true,
}

// rawTextElements contains tags whose text content must not be HTML-escaped
// because the browser parses them as raw text (JavaScript / CSS).
var rawTextElements = map[string]bool{
	"script": true,
	"style":  true,
}

// render writes the indented HTML for this element into sb.
func (e Element) render(sb *strings.Builder, indent string) {
	const indentStr = "  "

	// Nothing to render.
	if e.Tag == "" && e.TextContent == "" && len(e.ChildNodes) == 0 {
		return
	}

	// Tag-less node: either a raw fragment or a plain text/fragment node.
	if e.Tag == "" {
		if e.TextContent != "" {
			text := e.TextContent
			if !e.rawText {
				text = html.EscapeString(text)
			}
			sb.WriteString(indent + text + "\n")
		}
		for _, child := range e.ChildNodes {
			child.render(sb, indent+indentStr)
		}
		return
	}

	// Opening tag.
	sb.WriteString(indent + "<" + e.Tag)

	// Sort attribute keys so the output is deterministic.
	keys := make([]string, 0, len(e.Attributes))
	for k := range e.Attributes {
		if k != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := e.Attributes[k]
		sb.WriteString(" " + k)
		if v != "" {
			sb.WriteString(`="` + html.EscapeString(v) + `"`)
		}
	}

	// Void elements have no content or closing tag.
	if selfClosingTags[e.Tag] {
		sb.WriteString(" />\n")
		return
	}
	sb.WriteString(">\n")

	// Text content.
	if e.TextContent != "" {
		text := e.TextContent
		if !e.rawText && !rawTextElements[e.Tag] {
			text = html.EscapeString(text)
		}
		sb.WriteString(indent + indentStr + text + "\n")
	}

	// Child elements.
	for _, child := range e.ChildNodes {
		child.render(sb, indent+indentStr)
	}

	sb.WriteString(indent + "</" + e.Tag + ">\n")
}
