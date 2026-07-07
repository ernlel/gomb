package gomb_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ernlel/gomb"
)

// ── helpers ──────────────────────────────────────────────────────────────────

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

// ── basic rendering ──────────────────────────────────────────────────────────

func TestSimpleElement(t *testing.T) {
	out := gomb.E("p").T("Hello").ToString()
	assertContains(t, out, "<p>")
	assertContains(t, out, "Hello")
	assertContains(t, out, "</p>")
}

func TestAttributeRendered(t *testing.T) {
	out := gomb.E("div").A("class", "container").ToString()
	assertContains(t, out, `class="container"`)
}

func TestMultipleAttributes(t *testing.T) {
	out := gomb.E("a").A("href", "/home").A("class", "nav").ToString()
	assertContains(t, out, `href="/home"`)
	assertContains(t, out, `class="nav"`)
}

// Attribute order must be deterministic across multiple calls.
func TestAttributeOrder(t *testing.T) {
	out1 := gomb.E("div").A("z", "last").A("a", "first").A("m", "mid").ToString()
	out2 := gomb.E("div").A("z", "last").A("a", "first").A("m", "mid").ToString()
	if out1 != out2 {
		t.Errorf("attribute order not deterministic:\n%s\n%s", out1, out2)
	}
	// sorted: a, m, z
	idx_a := strings.Index(out1, `a="first"`)
	idx_m := strings.Index(out1, `m="mid"`)
	idx_z := strings.Index(out1, `z="last"`)
	if !(idx_a < idx_m && idx_m < idx_z) {
		t.Errorf("attributes not in alphabetical order: a=%d m=%d z=%d in:\n%s", idx_a, idx_m, idx_z, out1)
	}
}

func TestBooleanAttribute(t *testing.T) {
	out := gomb.E("input").A("type", "checkbox").A("checked", "").ToString()
	assertContains(t, out, "checked")
	assertNotContains(t, out, `checked="`)
}

func TestChildren(t *testing.T) {
	out := gomb.E("ul").C(
		gomb.E("li").T("one"),
		gomb.E("li").T("two"),
	).ToString()
	assertContains(t, out, "<ul>")
	assertContains(t, out, "<li>")
	assertContains(t, out, "one")
	assertContains(t, out, "two")
	assertContains(t, out, "</ul>")
}

// ── self-closing tags ────────────────────────────────────────────────────────

func TestSelfClosingTags(t *testing.T) {
	for _, tag := range []string{"br", "hr", "img", "input", "meta", "link"} {
		out := gomb.E(tag).ToString()
		assertContains(t, out, "/>")
		assertNotContains(t, out, "</"+tag+">")
	}
}

// ── HTML escaping ────────────────────────────────────────────────────────────

func TestTextEscaping(t *testing.T) {
	out := gomb.E("p").T(`<script>alert("xss")</script>`).ToString()
	assertNotContains(t, out, "<script>")
	assertContains(t, out, "&lt;script&gt;")
}

func TestAttributeEscaping(t *testing.T) {
	out := gomb.E("div").A("title", `say "hello" & goodbye`).ToString()
	assertNotContains(t, out, `"hello"`)
	assertContains(t, out, "&amp;")
}

// script and style content must NOT be entity-escaped.
func TestScriptContentNotEscaped(t *testing.T) {
	out := gomb.E("script").T(`if (a < b && c > d) alert("ok")`).ToString()
	assertContains(t, out, `if (a < b && c > d) alert("ok")`)
}

func TestStyleContentNotEscaped(t *testing.T) {
	out := gomb.E("style").T(`.btn::after { content: "»"; }`).ToString()
	assertContains(t, out, `.btn::after { content: "»"; }`)
}

// ── Raw ──────────────────────────────────────────────────────────────────────

func TestRawFragment(t *testing.T) {
	raw := `<span class="icon">★</span>`
	out := gomb.E("div").C(gomb.Raw(raw)).ToString()
	assertContains(t, out, raw)
}

func TestRawNotEscaped(t *testing.T) {
	out := gomb.E("div").C(gomb.Raw(`<b>bold</b>`)).ToString()
	assertNotContains(t, out, "&lt;")
}

// ── Fragment ─────────────────────────────────────────────────────────────────

func TestFragment(t *testing.T) {
	out := gomb.Fragment(
		gomb.E("li").T("a"),
		gomb.E("li").T("b"),
	).ToString()
	assertNotContains(t, out, "<>")  // no wrapper tag
	assertContains(t, out, "<li>")
}

// ── None / If / IfElse ───────────────────────────────────────────────────────

func TestNoneRendersNothing(t *testing.T) {
	out := gomb.None.ToString()
	if strings.TrimSpace(out) != "" {
		t.Errorf("None should render empty string, got %q", out)
	}
}

func TestIfTrue(t *testing.T) {
	out := gomb.E("div").C(gomb.If(true, gomb.E("span").T("yes"))).ToString()
	assertContains(t, out, "yes")
}

func TestIfFalse(t *testing.T) {
	out := gomb.E("div").C(gomb.If(false, gomb.E("span").T("yes"))).ToString()
	assertNotContains(t, out, "yes")
}

func TestIfElseTrue(t *testing.T) {
	out := gomb.IfElse(true, gomb.E("span").T("yes"), gomb.E("span").T("no")).ToString()
	assertContains(t, out, "yes")
	assertNotContains(t, out, "no")
}

func TestIfElseFalse(t *testing.T) {
	out := gomb.IfElse(false, gomb.E("span").T("yes"), gomb.E("span").T("no")).ToString()
	assertNotContains(t, out, "yes")
	assertContains(t, out, "no")
}

// ── Map ──────────────────────────────────────────────────────────────────────

func TestMap(t *testing.T) {
	items := []string{"Alice", "Bob", "Carol"}
	out := gomb.E("ul").C(gomb.Map(items, func(s string) gomb.Element {
		return gomb.E("li").T(s)
	})...).ToString()
	for _, name := range items {
		assertContains(t, out, name)
	}
}

func TestMapEmpty(t *testing.T) {
	out := gomb.E("ul").C(gomb.Map([]string{}, func(s string) gomb.Element {
		return gomb.E("li").T(s)
	})...).ToString()
	assertContains(t, out, "<ul>")
	assertNotContains(t, out, "<li>")
}

// ── Render (io.Writer) ───────────────────────────────────────────────────────

func TestRender(t *testing.T) {
	var buf bytes.Buffer
	err := gomb.E("p").T("hi").Render(&buf)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	assertContains(t, buf.String(), "<p>")
	assertContains(t, buf.String(), "hi")
}

// ── immutability (no shared-map mutation) ────────────────────────────────────

func TestAttributeImmutability(t *testing.T) {
	base := gomb.E("div").A("class", "base")
	derived := base.A("id", "unique")

	baseOut := base.ToString()
	if strings.Contains(baseOut, `id="unique"`) {
		t.Error("mutating derived element corrupted base element attributes")
	}
	derivedOut := derived.ToString()
	assertContains(t, derivedOut, `id="unique"`)
}

func TestChildrenImmutability(t *testing.T) {
	base := gomb.E("ul").C(gomb.E("li").T("first"))
	_ = base.C(gomb.E("li").T("second"))

	baseOut := base.ToString()
	if strings.Contains(baseOut, "second") {
		t.Error("mutating derived element corrupted base element children")
	}
}

// ── nested structure ─────────────────────────────────────────────────────────

func TestNestedElements(t *testing.T) {
	out := gomb.E("html").A("lang", "en").C(
		gomb.E("head").C(gomb.E("title").T("Test")),
		gomb.E("body").C(
			gomb.E("h1").T("Hello"),
			gomb.E("p").T("World"),
		),
	).ToString()
	assertContains(t, out, `lang="en"`)
	assertContains(t, out, "<title>")
	assertContains(t, out, "Test")
	assertContains(t, out, "<h1>")
	assertContains(t, out, "Hello")
}

// ── Render nil writer ────────────────────────────────────────────────────────

func TestRenderNilWriter(t *testing.T) {
	err := gomb.E("p").T("hi").Render(nil)
	if err == nil {
		t.Error("Render with nil writer should return an error")
	}
}

// ── Classes ──────────────────────────────────────────────────────────────────

func TestClassesBasic(t *testing.T) {
	s := gomb.Classes("foo", "bar", "baz")
	if s != "foo bar baz" {
		t.Errorf("expected 'foo bar baz', got %q", s)
	}
}

func TestClassesSkipsEmpty(t *testing.T) {
	s := gomb.Classes("foo", "", "bar", "", "baz")
	if s != "foo bar baz" {
		t.Errorf("expected 'foo bar baz', got %q", s)
	}
}

func TestClassesSingleArg(t *testing.T) {
	s := gomb.Classes("only")
	if s != "only" {
		t.Errorf("expected 'only', got %q", s)
	}
}

func TestClassesAllEmpty(t *testing.T) {
	s := gomb.Classes("", "", "")
	if s != "" {
		t.Errorf("expected empty string, got %q", s)
	}
}

// ── Data ─────────────────────────────────────────────────────────────────────

func TestData(t *testing.T) {
	out := gomb.E("div").Data("count", "0").Data("user", "42").ToString()
	assertContains(t, out, `data-count="0"`)
	assertContains(t, out, `data-user="42"`)
}

// ── Style ────────────────────────────────────────────────────────────────────

func TestStyle(t *testing.T) {
	out := gomb.E("div").Style("color: red; font-weight: bold").ToString()
	assertContains(t, out, `style="color: red; font-weight: bold"`)
}

// ── When ─────────────────────────────────────────────────────────────────────

func TestWhenTrue(t *testing.T) {
	called := false
	out := gomb.E("div").C(gomb.When(true, func() gomb.Element {
		called = true
		return gomb.E("span").T("yes")
	})).ToString()
	if !called {
		t.Error("When with true condition should call fn")
	}
	assertContains(t, out, "yes")
}

func TestWhenFalse(t *testing.T) {
	called := false
	out := gomb.E("div").C(gomb.When(false, func() gomb.Element {
		called = true
		return gomb.E("span").T("yes")
	})).ToString()
	if called {
		t.Error("When with false condition should NOT call fn")
	}
	assertNotContains(t, out, "yes")
}

func TestAttributeOverwrite(t *testing.T) {
	out := gomb.E("div").A("class", "a").A("class", "b").ToString()
	if strings.Count(out, `class="`) != 1 {
		t.Errorf("overwrite should produce exactly one class attr, got: %q", out)
	}
	if !strings.Contains(out, `class="b"`) {
		t.Errorf("last write should win, got: %q", out)
	}
}

func TestTextAndChildren(t *testing.T) {
	out := gomb.E("div").T("before").C(gomb.E("span").T("child")).ToString()
	assertContains(t, out, "before")
	assertContains(t, out, "child")
}

func TestCNoArgs(t *testing.T) {
	out := gomb.E("div").C().ToString()
	assertContains(t, out, "<div>")
	assertContains(t, out, "</div>")
}

func TestRawEmpty(t *testing.T) {
	out := gomb.Raw("").ToString()
	if strings.TrimSpace(out) != "" {
		t.Errorf("empty Raw should render nothing, got: %q", out)
	}
}

// ── Attrs / NS / With ────────────────────────────────────────────────────────

func TestAttrs(t *testing.T) {
	out := gomb.E("div").Attrs(
		gomb.Attr{"class", "container"},
		gomb.Attr{"id", "main"},
	).ToString()
	assertContains(t, out, `class="container"`)
	assertContains(t, out, `id="main"`)
}

func TestAttrsEmpty(t *testing.T) {
	out := gomb.E("div").Attrs().ToString()
	assertContains(t, out, "<div>")
	assertContains(t, out, "</div>")
}

func TestNS(t *testing.T) {
	hx := gomb.NS("hx-")
	out := gomb.E("button").Attrs(
		hx("swap", "outerHTML"),
		hx("target", "#result"),
		hx("get", "/api/data"),
	).ToString()
	assertContains(t, out, `hx-swap="outerHTML"`)
	assertContains(t, out, `hx-target="#result"`)
	assertContains(t, out, `hx-get="/api/data"`)
}

func TestNSDataPrefix(t *testing.T) {
	data := gomb.NS("data-")
	out := gomb.E("li").Attrs(
		data("user", "42"),
		data("role", "admin"),
	).ToString()
	assertContains(t, out, `data-user="42"`)
	assertContains(t, out, `data-role="admin"`)
}

func TestNSMixedWithA(t *testing.T) {
	hx := gomb.NS("hx-")
	out := gomb.E("button").
		A("class", "btn").
		Attrs(hx("swap", "outerHTML"), hx("target", "#result")).
		ToString()
	assertContains(t, out, `class="btn"`)
	assertContains(t, out, `hx-swap="outerHTML"`)
	assertContains(t, out, `hx-target="#result"`)
}

func TestWith(t *testing.T) {
	btn := func(e gomb.Element) gomb.Element {
		return e.A("type", "submit").A("class", "btn")
	}
	out := gomb.E("button").T("Save").With(btn).ToString()
	assertContains(t, out, `type="submit"`)
	assertContains(t, out, `class="btn"`)
	assertContains(t, out, "Save")
}

func TestWithMultiple(t *testing.T) {
	setType := func(e gomb.Element) gomb.Element { return e.A("type", "text") }
	addClasses := func(e gomb.Element) gomb.Element { return e.A("class", "input") }
	out := gomb.E("input").With(setType, addClasses).ToString()
	assertContains(t, out, `type="text"`)
	assertContains(t, out, `class="input"`)
}

func TestWithNoArgs(t *testing.T) {
	out := gomb.E("p").T("hi").With().ToString()
	assertContains(t, out, "hi")
}

// ── aliases ──────────────────────────────────────────────────────────────────

func TestElAlias(t *testing.T) {
	a := gomb.E("div").A("class", "test").T("hello").ToString()
	b := gomb.El("div").A("class", "test").T("hello").ToString()
	if a != b {
		t.Errorf("E and El should produce identical output:\nE:\n%s\nEl:\n%s", a, b)
	}
}

func TestAttrAlias(t *testing.T) {
	a := gomb.E("div").A("class", "test").ToString()
	b := gomb.E("div").Attr("class", "test").ToString()
	if a != b {
		t.Errorf("A and Attr should produce identical output:\nA:\n%s\nAttr:\n%s", a, b)
	}
}

func TestTextAlias(t *testing.T) {
	a := gomb.E("p").T("hello").ToString()
	b := gomb.E("p").Text("hello").ToString()
	if a != b {
		t.Errorf("T and Text should produce identical output:\nT:\n%s\nText:\n%s", a, b)
	}
}

func TestChildrenAlias(t *testing.T) {
	a := gomb.E("ul").C(gomb.E("li").T("one")).ToString()
	b := gomb.E("ul").Children(gomb.E("li").T("one")).ToString()
	if a != b {
		t.Errorf("C and Children should produce identical output:\nC:\n%s\nChildren:\n%s", a, b)
	}
}

func TestAsAlias(t *testing.T) {
	a := gomb.E("div").Attrs(gomb.Attr{"class", "a"}, gomb.Attr{"id", "b"}).ToString()
	b := gomb.E("div").As(gomb.Attr{"class", "a"}, gomb.Attr{"id", "b"}).ToString()
	if a != b {
		t.Errorf("Attrs and As should produce identical output:\nAttrs:\n%s\nAs:\n%s", a, b)
	}
}

// ── Txt ──────────────────────────────────────────────────────────────────────

func TestTxt(t *testing.T) {
	out := gomb.Txt("hello").ToString()
	assertContains(t, out, "hello")
}

func TestTxtEscaped(t *testing.T) {
	out := gomb.Txt(`<b>`).ToString()
	assertNotContains(t, out, "<b>")
	assertContains(t, out, "&lt;b&gt;")
}
