package gomb

import (
	"strings"
)

type Element struct {
	tag      string
	attrs    map[string]string
	children []Element
	text     string
}

func E(tag string) Element {
	return Element{tag: tag}
}

func (e Element) A(key, value string) Element {
	if e.attrs == nil {
		e.attrs = make(map[string]string)
	}
	e.attrs[key] = value
	return e
}

func (e Element) T(value string) Element {
	e.text = value
	return e
}

func (e Element) C(elements ...Element) Element {
	e.children = append(e.children, elements...)
	return e
}

func (e Element) ToString() string {
	return e.toStringIndented("")
}

var selfClosingTags = []string{
	"area", "base", "br", "col", "embed", "hr",
	"img", "input", "link", "meta", "param", "source",
	"track", "wbr",
}

func isSelfClosing(tag string) bool {
	for _, t := range selfClosingTags {
		if t == tag {
			return true
		}
	}
	return false
}

func (e Element) toStringIndented(indent string) string {
	const indentString = "  "
	var sb strings.Builder

	if e.tag != "" {
		sb.WriteString(indent + "<" + e.tag)
		for k, v := range e.attrs {
			if k == "" {
				continue
			}
			sb.WriteString(" ")
			sb.WriteString(k)
			if v != "" {
				sb.WriteString(`="`)
				sb.WriteString(v)
				sb.WriteString(`"`)
			}
		}

		if isSelfClosing(e.tag) {
			sb.WriteString(" />\n")
			return sb.String()
		} else {
			sb.WriteString(">\n")
		}
	}

	if e.text != "" {
		sb.WriteString(indent + indentString + e.text + "\n")
	}

	for _, child := range e.children {
		sb.WriteString(child.toStringIndented(indent + indentString))
	}

	if e.tag != "" {
		sb.WriteString(indent + "</" + e.tag + ">\n")
	}

	return sb.String()
}
