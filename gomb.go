// Package gomb provides a fluent, type-safe API for building HTML programmatically in Go.
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
	"errors"
	"html"
	"io"
	"sort"
	"strings"
)

// ErrNilWriter is returned by Render when called with a nil io.Writer.
var ErrNilWriter = errors.New("gomb: Render called with nil writer")

// Element represents an HTML element.
type Element struct {
	Tag        string
	Attributes map[string]string
	ChildNodes []*Element
	TextContent string
	rawText    bool
}

// None is a nil *Element. It renders to nothing and is the canonical return value
// when a conditional helper has no output.
var None *Element

// E creates a new HTML element with the given tag name.
func E(tag string) *Element {
	return &Element{Tag: tag}
}

// El is an alias for E.
func El(tag string) *Element {
	return E(tag)
}

// Txt creates a tag-less text element.
func Txt(content string) *Element {
	return &Element{TextContent: content}
}

// Raw creates an element whose content is written verbatim without HTML escaping.
func Raw(rawHTML string) *Element {
	return &Element{TextContent: rawHTML, rawText: true}
}

// Fragment creates a tag-less wrapper around multiple elements.
func Fragment(elements ...*Element) *Element {
	return E("").C(elements...)
}

// A sets attributes from key-value pairs. Returns the element for chaining.
// Odd trailing arguments are silently ignored.
func (e *Element) A(pairs ...string) *Element {
	if e.Attributes == nil {
		e.Attributes = make(map[string]string)
	}
	for i := 0; i+1 < len(pairs); i += 2 {
		e.Attributes[pairs[i]] = pairs[i+1]
	}
	return e
}

// Attr is an alias for A.
func (e *Element) Attr(pairs ...string) *Element {
	return e.A(pairs...)
}

// T sets text content. Returns the element for chaining.
func (e *Element) T(value string) *Element {
	e.TextContent = value
	e.rawText = false
	return e
}

// Text is an alias for T.
func (e *Element) Text(value string) *Element {
	return e.T(value)
}

// C appends child elements. Returns the element for chaining.
func (e *Element) C(elements ...*Element) *Element {
	for _, el := range elements {
		if el != nil {
			e.ChildNodes = append(e.ChildNodes, el)
		}
	}
	return e
}

// Children is an alias for C.
func (e *Element) Children(elements ...*Element) *Element {
	return e.C(elements...)
}

// ToString renders the element tree to an indented HTML string.
func (e *Element) ToString() string {
	if e == nil {
		return ""
	}
	var sb strings.Builder
	e.render(&sb, "")
	return sb.String()
}

// Render writes the HTML representation of the element to w.
// Returns the number of bytes written and any error encountered.
func (e *Element) Render(w io.Writer) (int64, error) {
	if e == nil {
		return 0, nil
	}
	if w == nil {
		return 0, ErrNilWriter
	}
	n, err := io.WriteString(w, e.ToString())
	return int64(n), err
}

// Classes returns a space-joined string from the given class names, skipping empty strings.
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

// Data is a shorthand for data-* attributes.
func (e *Element) Data(key, value string) *Element {
	return e.A("data-"+key, value)
}

// Style is a shorthand for the style attribute.
func (e *Element) Style(css string) *Element {
	return e.A("style", css)
}

// Class is a shorthand for setting the class attribute.
// It joins the given names using Classes() and sets "class".
func (e *Element) Class(names ...string) *Element {
	return e.A("class", Classes(names...))
}

// Id is a shorthand for setting the id attribute.
func (e *Element) Id(id string) *Element {
	return e.A("id", id)
}

// When applies fn to the element when cond is true. Returns the element for chaining.
func (e *Element) When(cond bool, fn func(*Element)) *Element {
	if cond {
		fn(e)
	}
	return e
}

// Clone returns a shallow copy of the element.
// The Attributes map is copied, but ChildNodes slice references are shared.
func (e *Element) Clone() *Element {
	if e == nil {
		return nil
	}
	clone := &Element{
		Tag:         e.Tag,
		TextContent: e.TextContent,
		rawText:     e.rawText,
	}
	if e.Attributes != nil {
		clone.Attributes = make(map[string]string, len(e.Attributes))
		for k, v := range e.Attributes {
			clone.Attributes[k] = v
		}
	}
	if len(e.ChildNodes) > 0 {
		clone.ChildNodes = make([]*Element, len(e.ChildNodes))
		copy(clone.ChildNodes, e.ChildNodes)
	}
	return clone
}

// Attr is a key-value attribute pair.
type Attr struct {
	Key, Value string
}

// Attrs applies multiple Attr pairs to the element.
func (e *Element) Attrs(attrs ...Attr) *Element {
	for _, a := range attrs {
		e.A(a.Key, a.Value)
	}
	return e
}

// As is an alias for Attrs.
func (e *Element) As(attrs ...Attr) *Element {
	return e.Attrs(attrs...)
}

// NS returns a function that builds namespaced attribute pairs.
func NS(prefix string) func(key, value string) Attr {
	return func(key, value string) Attr {
		return Attr{Key: prefix + key, Value: value}
	}
}

// With applies the given transformer functions to the element.
func (e *Element) With(fns ...func(*Element)) *Element {
	for _, fn := range fns {
		fn(e)
	}
	return e
}

// When returns el from fn() when cond is true, otherwise returns None.
func When(cond bool, fn func() *Element) *Element {
	if cond {
		return fn()
	}
	return nil
}

// If returns el when cond is true, otherwise returns None.
func If(cond bool, el *Element) *Element {
	if cond {
		return el
	}
	return nil
}

// IfElse returns ifEl when cond is true, otherwise returns elseEl.
func IfElse(cond bool, ifEl, elseEl *Element) *Element {
	if cond {
		return ifEl
	}
	return elseEl
}

// Map transforms a slice of values into a slice of Elements.
func Map[T any](items []T, fn func(T) *Element) []*Element {
	result := make([]*Element, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}

var selfClosingTags = map[string]bool{
	"area": true, "base": true, "br": true, "col": true,
	"embed": true, "hr": true, "img": true, "input": true,
	"link": true, "meta": true, "param": true, "source": true,
	"track": true, "wbr": true,
}

var rawTextElements = map[string]bool{
	"script": true,
	"style":  true,
}

func (e *Element) render(sb *strings.Builder, indent string) {
	const indentStr = "  "

	if e.Tag == "" && e.TextContent == "" && len(e.ChildNodes) == 0 {
		return
	}

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

	sb.WriteString(indent + "<" + e.Tag)

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

	if selfClosingTags[e.Tag] {
		sb.WriteString(" />\n")
		return
	}
	sb.WriteString(">\n")

	if e.TextContent != "" {
		text := e.TextContent
		if !e.rawText && !rawTextElements[e.Tag] {
			text = html.EscapeString(text)
		}
		sb.WriteString(indent + indentStr + text + "\n")
	}

	for _, child := range e.ChildNodes {
		child.render(sb, indent+indentStr)
	}

	sb.WriteString(indent + "</" + e.Tag + ">\n")
}
