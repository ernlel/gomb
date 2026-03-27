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
type Element struct {
	Tag      string
	Attrs    map[string]string
	Children []Element
	Text     string
	rawText  bool // when true, Text is written verbatim (no HTML escaping)
}

// None is the zero Element. It renders to nothing and is the canonical return value
// when a conditional helper has no output.
var None = Element{}

// E creates a new HTML element with the given tag name.
// Use an empty tag ("") to create a tag-less fragment.
func E(tag string) Element {
	return Element{Tag: tag}
}

// Raw creates an element whose content is written verbatim to the output without any
// HTML escaping. Use it for pre-rendered HTML fragments or for content inside
// <script> and <style> tags that must not be entity-encoded.
//
//	E("script").A("type", "text/javascript").C(Raw(`console.log("hello")`))
//	E("style").C(Raw(`.btn { color: red; }`))
//	E("div").C(Raw("<span>pre-rendered</span>"))
func Raw(rawHTML string) Element {
	return Element{Text: rawHTML, rawText: true}
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
	newAttrs := make(map[string]string, len(e.Attrs)+1)
	for k, v := range e.Attrs {
		newAttrs[k] = v
	}
	newAttrs[key] = value
	e.Attrs = newAttrs
	return e
}

// T returns a new Element with the given text content. The text is HTML-escaped
// when the element is rendered, preventing XSS when displaying user-supplied data.
// For unescaped output use Raw().
func (e Element) T(value string) Element {
	e.Text = value
	e.rawText = false
	return e
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
	newChildren := make([]Element, len(e.Children)+len(elements))
	copy(newChildren, e.Children)
	copy(newChildren[len(e.Children):], elements)
	e.Children = newChildren
	return e
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
	_, err := io.WriteString(w, e.ToString())
	return err
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
	if e.Tag == "" && e.Text == "" && len(e.Children) == 0 {
		return
	}

	// Tag-less node: either a raw fragment or a plain text/fragment node.
	if e.Tag == "" {
		if e.Text != "" {
			text := e.Text
			if !e.rawText {
				text = html.EscapeString(text)
			}
			sb.WriteString(indent + indentStr + text + "\n")
		}
		for _, child := range e.Children {
			child.render(sb, indent+indentStr)
		}
		return
	}

	// Opening tag.
	sb.WriteString(indent + "<" + e.Tag)

	// Sort attribute keys so the output is deterministic.
	keys := make([]string, 0, len(e.Attrs))
	for k := range e.Attrs {
		if k != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := e.Attrs[k]
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
	if e.Text != "" {
		text := e.Text
		if !e.rawText && !rawTextElements[e.Tag] {
			text = html.EscapeString(text)
		}
		sb.WriteString(indent + indentStr + text + "\n")
	}

	// Child elements.
	for _, child := range e.Children {
		child.render(sb, indent+indentStr)
	}

	sb.WriteString(indent + "</" + e.Tag + ">\n")
}
