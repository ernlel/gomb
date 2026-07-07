// Example: named HTML element constructors with inline and chained styles.
//
// Run:
//
//	go run main.go
//	open http://localhost:8080
package main

import (
	"fmt"
	"net/http"
	"strings"

	. "github.com/ernlel/gomb"
	. "github.com/ernlel/gomb/pkg/html"
)

// ── components ───────────────────────────────────────────────────────────────

func NavItem(label, href string, active bool) Element {
	cls := "px-3 py-2 rounded text-sm font-medium"
	if active {
		cls += " bg-gray-900 text-white"
	} else {
		cls += " text-gray-300 hover:bg-gray-700"
	}
	return Li(
		Anchor().Attr("href", href).Attr("class", cls).Text(label),
	)
}

func Card(title, body string) Element {
	return Div(
		Attr{Key: "class", Value: "bg-white rounded-lg shadow p-6"},
		H3(Txt(title), Attr{Key: "class", Value: "text-lg font-semibold mb-2"}),
		P(Txt(body), Attr{Key: "class", Value: "text-gray-600"}),
	)
}

func StatCard(label, value string) Element {
	return Div(
		Attr{Key: "class", Value: "bg-white rounded-lg shadow p-4 text-center"},
		Div(Txt(value), Attr{Key: "class", Value: "text-3xl font-bold text-indigo-600"}),
		Div(Txt(label), Attr{Key: "class", Value: "text-sm text-gray-500 mt-1"}),
	)
}

// ── page ─────────────────────────────────────────────────────────────────────

func page() Element {
	return Html(
		Attr{Key: "lang", Value: "en"},
		Head(
			Meta(Attr{Key: "charset", Value: "UTF-8"}),
			Meta(Attr{Key: "name", Value: "viewport"}, Attr{Key: "content", Value: "width=device-width, initial-scale=1"}),
			TitleElement(Txt("gomb — Named Constructors")),
			ScriptElement(Txt(""), Attr{Key: "src", Value: "https://cdn.tailwindcss.com"}),
		),
		Body(
			Attr{Key: "class", Value: "bg-gray-100 min-h-screen"},

			// ── nav ──
			Nav(
				Attr{Key: "class", Value: "bg-gray-800"},
				Div(
					Attr{Key: "class", Value: "max-w-7xl mx-auto px-4"},
					Div(
						Attr{Key: "class", Value: "flex items-center justify-between h-16"},
						Div(
							Attr{Key: "class", Value: "flex items-center"},
							Div(Txt("gomb"), Attr{Key: "class", Value: "text-white font-bold text-lg"}),
						),
						Div(
							Attr{Key: "class", Value: "flex space-x-4"},
							NavItem("Home", "/", true),
							NavItem("Docs", "/docs", false),
							NavItem("Examples", "/examples", false),
						),
					),
				),
			),

			// ── hero ──
			Div(
				Attr{Key: "class", Value: "max-w-7xl mx-auto py-12 px-4"},
				Div(
					Attr{Key: "class", Value: "text-center mb-12"},
					H1(
						Txt("Build HTML with Go"),
						Attr{Key: "class", Value: "text-4xl font-bold text-gray-900 mb-4"},
					),
					P(
						Txt("Type-safe. Auto-escaped. No templates."),
						Attr{Key: "class", Value: "text-xl text-gray-600"},
					),
				),

				// ── stats ──
				Div(
					Attr{Key: "class", Value: "grid grid-cols-3 gap-6 mb-12"},
					StatCard("Elements", "114"),
					StatCard("Test Coverage", "99%"),
					StatCard("Dependencies", "0"),
				),

				// ── cards ──
				H2(
					Txt("Why gomb?"),
					Attr{Key: "class", Value: "text-2xl font-bold text-gray-900 mb-6"},
				),
				Div(
					Attr{Key: "class", Value: "grid grid-cols-2 gap-6"},
					Card("Type Safety",
						"The Go compiler catches typos — no runtime template errors."),
					Card("Auto-Escaping",
						"All text and attributes are HTML-escaped. XSS is off by default."),
					Card("IDE Support",
						"Full autocomplete, rename refactoring, and GoDoc on every element."),
					Card("Zero Dependencies",
						"The core library imports only the standard library."),
				),

				// ── code sample ──
				H2(
					Txt("Inline style — looks like HTML"),
					Attr{Key: "class", Value: "text-2xl font-bold text-gray-900 mt-12 mb-4"},
				),
				Pre(
					Code(Txt(strings.TrimSpace(`
page := Html(
    Head(
        Meta(Attr{Key:"charset", Value:"UTF-8"}),
        TitleElement(Txt("My App")),
    ),
    Body(
        H1(Txt("Hello")),
        P(Txt("Welcome to gomb.")),
    ),
).Attr("lang","en")`[1:])),
						Attr{Key: "class", Value: "text-sm"},
					),
					Attr{Key: "class", Value: "bg-gray-900 text-gray-100 rounded-lg p-6 overflow-x-auto"},
				),

				// ── footer ──
				Footer(
					Attr{Key: "class", Value: "mt-16 text-center text-sm text-gray-500"},
					P(Txt("Built with gomb — Go Markup Builder")),
				),
			),
		),
	)
}

// ── server ───────────────────────────────────────────────────────────────────

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<!DOCTYPE html>\n")
	page().Render(w)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
