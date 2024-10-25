package gomb

import (
	"strings"
)

type Element struct {
	Tag      string
	Attrs    map[string]string
	Children []Element
	Text     string
}

func E(tag string) Element {
	return Element{Tag: tag}
}

func (e Element) A(key, value string) Element {
	if e.Attrs == nil {
		e.Attrs = make(map[string]string)
	}
	e.Attrs[key] = value
	return e
}

func (e Element) T(value string) Element {
	e.Text = value
	return e
}

func (e Element) C(elements ...Element) Element {
	e.Children = append(e.Children, elements...)
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

	if e.Tag != "" {
		sb.WriteString(indent + "<" + e.Tag)
		for k, v := range e.Attrs {
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

		if isSelfClosing(e.Tag) {
			sb.WriteString(" />\n")
			return sb.String()
		} else {
			sb.WriteString(">\n")
		}
	}

	if e.Text != "" {
		sb.WriteString(indent + indentString + e.Text + "\n")
	}

	for _, child := range e.Children {
		sb.WriteString(child.toStringIndented(indent + indentString))
	}

	if e.Tag != "" {
		sb.WriteString(indent + "</" + e.Tag + ">\n")
	}

	return sb.String()
}
