package html_test

import (
	"strings"
	"testing"

	"github.com/ernlel/gomb"
	. "github.com/ernlel/gomb/pkg/html"
)

func assertContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Errorf("output does not contain %q\ngot:\n%s", want, got)
	}
}

func assertNotContains(t *testing.T, got, want string) {
	t.Helper()
	if strings.Contains(got, want) {
		t.Errorf("output should not contain %q\ngot:\n%s", want, got)
	}
}

func TestInlineMatchesChained(t *testing.T) {
	a := Div().Attr("class", "card").Children(
		H2().Text("Title"),
		P().Text("Body"),
	).ToString()

	b := Div(
		gomb.Attr{Key: "class", Value: "card"},
		H2(gomb.Txt("Title")),
		P(gomb.Txt("Body")),
	).ToString()

	if a != b {
		t.Errorf("chained and inline should produce identical output:\nchained:\n%s\ninline:\n%s", a, b)
	}
}

func TestInlineMixed(t *testing.T) {
	out := Div(
		gomb.Attr{Key: "class", Value: "container"},
		gomb.Attr{Key: "id", Value: "main"},
		Span(gomb.Txt("hello")),
	).ToString()
	assertContains(t, out, `class="container"`)
	assertContains(t, out, `id="main"`)
	assertContains(t, out, "hello")
}

func TestAllElementsExist(t *testing.T) {
	elements := map[string]func(args ...interface{}) *gomb.Element{
		"a":       func(args ...interface{}) *gomb.Element { return Anchor(args...) },
		"div":     func(args ...interface{}) *gomb.Element { return Div(args...) },
		"span":    func(args ...interface{}) *gomb.Element { return Span(args...) },
		"p":       func(args ...interface{}) *gomb.Element { return P(args...) },
		"h1":      func(args ...interface{}) *gomb.Element { return H1(args...) },
		"html":    func(args ...interface{}) *gomb.Element { return Html(args...) },
		"head":    func(args ...interface{}) *gomb.Element { return Head(args...) },
		"body":    func(args ...interface{}) *gomb.Element { return Body(args...) },
		"meta":    func(args ...interface{}) *gomb.Element { return Meta(args...) },
		"input":   func(args ...interface{}) *gomb.Element { return InputElement(args...) },
		"script":  func(args ...interface{}) *gomb.Element { return ScriptElement(args...) },
		"style":   func(args ...interface{}) *gomb.Element { return StyleElement(args...) },
		"title":   func(args ...interface{}) *gomb.Element { return TitleElement(args...) },
		"var":     func(args ...interface{}) *gomb.Element { return VarElement(args...) },
		"map":     func(args ...interface{}) *gomb.Element { return MapElement(args...) },
		"select":  func(args ...interface{}) *gomb.Element { return SelectElement(args...) },
		"data":    func(args ...interface{}) *gomb.Element { return DataElement(args...) },
		"template": func(args ...interface{}) *gomb.Element { return TemplateElement(args...) },
	}

	for tag, fn := range elements {
		out := fn().ToString()
		if !strings.Contains(out, "<"+tag) {
			t.Errorf("element <%s>: output missing opening tag, got:\n%s", tag, out)
		}
	}
}
