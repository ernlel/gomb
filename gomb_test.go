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
