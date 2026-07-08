// Package main demonstrates using gomb with Alpine.js for client-side interactivity.
//
// Alpine.js directives (x-data, x-show, x-model, @click, …) are just HTML
// attributes and map directly to gomb's .A() method. Because attribute values are
// HTML-escaped, JSON objects used in x-data are written as-is and arrive at the
// browser correctly decoded (&#34; → ").
//
// Run:
//
//	go run main.go
//	open http://localhost:8080
package main

import (
	"fmt"
	"net/http"

	. "github.com/ernlel/gomb"
)

// ── components ───────────────────────────────────────────────────────────────

// accordion renders a collapsible section using Alpine.js x-show / @click.
func accordion(title, body string) *Element {
	return E("div").
		A("x-data", `{"open": false}`).
		A("class", "border rounded mb-2").
		C(
			E("button").
				A("class", "w-full text-left px-4 py-3 font-semibold flex justify-between items-center").
				A("@click", "open = !open").
				C(
					E("span").T(title),
					E("span").A("x-text", `open ? "▲" : "▼"`),
				),
			E("div").
				A("x-show", "open").
				A("class", "px-4 py-3 text-gray-600 border-t").
				T(body),
		)
}

// tabPanel renders a tabbed interface.
func tabPanel() *Element {
	tabs := []struct{ label, content string }{
		{"Overview", "This is the overview tab."},
		{"Details", "Here are more details about the product."},
		{"Reviews", "Customers love this product – 4.8 ★"},
	}

	tabButtons := Map(tabs, func(t struct{ label, content string }) *Element {
		return E("button").
			A("class", "px-4 py-2 text-sm font-medium").
			A(":class", fmt.Sprintf(`active === %q ? "border-b-2 border-blue-600 text-blue-600" : "text-gray-500 hover:text-gray-700"`, t.label)).
			A("@click", fmt.Sprintf(`active = %q`, t.label)).
			T(t.label)
	})

	tabContents := Map(tabs, func(t struct{ label, content string }) *Element {
		return E("div").
			A("x-show", fmt.Sprintf(`active === %q`, t.label)).
			A("class", "p-4 text-gray-700").
			T(t.content)
	})

	return E("div").
		A("x-data", `{"active": "Overview"}`).
		A("class", "border rounded").
		C(
			E("div").A("class", "flex border-b").C(tabButtons...),
			E("div").C(tabContents...),
		)
}

// modal renders a dialog that can be opened/closed via Alpine state.
func modal() *Element {
	return E("div").
		A("x-data", `{"open": false}`).
		C(
			E("button").
				A("class", "bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700").
				A("@click", "open = true").
				T("Open Modal"),
			// Backdrop
			E("div").
				A("x-show", "open").
				A("class", "fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50").
				A("@click.self", "open = false").
				C(
					E("div").
						A("class", "bg-white rounded shadow-lg p-6 max-w-sm w-full").
						C(
							E("h2").A("class", "text-xl font-bold mb-2").T("Hello from Alpine!"),
							E("p").A("class", "text-gray-600 mb-4").T("This modal is powered by gomb + Alpine.js."),
							E("button").
								A("class", "bg-gray-200 px-4 py-2 rounded hover:bg-gray-300").
								A("@click", "open = false").
								T("Close"),
						),
				),
		)
}

// counter renders a live counter using x-data and x-text.
func counter() *Element {
	return E("div").
		A("x-data", `{"count": 0}`).
		A("class", "flex items-center gap-4").
		C(
			E("button").
				A("class", "bg-red-500 text-white w-8 h-8 rounded-full font-bold").
				A("@click", "count--").
				T("-"),
			E("span").
				A("class", "text-2xl font-mono w-12 text-center").
				A("x-text", "count"),
			E("button").
				A("class", "bg-green-500 text-white w-8 h-8 rounded-full font-bold").
				A("@click", "count++").
				T("+"),
		)
}

// liveSearch renders a list filtered as the user types (no server round-trip).
func liveSearch(items []string) *Element {
	// Encode the items slice as a JS array literal for use in x-data.
	quotedItems := make([]string, len(items))
	for i, item := range items {
		quotedItems[i] = fmt.Sprintf("%q", item)
	}
	jsArray := "["
	for i, q := range quotedItems {
		if i > 0 {
			jsArray += ", "
		}
		jsArray += q
	}
	jsArray += "]"

	return E("div").
		A("x-data", fmt.Sprintf(`{"query": "", "items": %s}`, jsArray)).
		A("class", "space-y-2").
		C(
			E("input").
				A("type", "search").
				A("placeholder", "Filter…").
				A("x-model", "query").
				A("class", "w-full border rounded px-3 py-2 text-sm"),
			E("ul").
				A("class", "border rounded divide-y").
				C(
					E("template").
						A("x-for", "item in items.filter(i => i.toLowerCase().includes(query.toLowerCase()))").
						C(
							E("li").
								A("class", "px-3 py-2 text-sm").
								A("x-text", "item"),
						),
				),
		)
}

// page assembles the full demo page.
func page() *Element {
	languages := []string{"Go", "TypeScript", "Rust", "Python", "Kotlin", "Swift", "Elixir", "Zig"}

	return E("html").A("lang", "en").C(
		E("head").C(
			E("meta").A("charset", "UTF-8"),
			E("meta").A("name", "viewport").A("content", "width=device-width, initial-scale=1"),
			E("title").T("gomb + Alpine.js"),
			E("link").
				A("rel", "stylesheet").
				A("href", "https://cdn.jsdelivr.net/npm/tailwindcss@3/dist/tailwind.min.css"),
			// Alpine.js – loaded as a module so x-cloak hides elements until ready.
			// In production pin the version and add integrity+crossorigin attributes.
			// Generate the hash at https://www.srihash.org/
			E("script").
				A("src", "https://unpkg.com/alpinejs@3").
				A("crossorigin", "anonymous").
				A("defer", ""),
			E("style").T(`[x-cloak] { display: none !important; }`),
		),
		E("body").A("class", "bg-gray-50 min-h-screen p-8 space-y-10").A("x-cloak", "").C(
			E("div").A("class", "max-w-lg mx-auto space-y-10").C(
				E("h1").A("class", "text-3xl font-bold").T("gomb + Alpine.js"),

				E("section").C(
					E("h2").A("class", "text-xl font-semibold mb-3").T("Counter"),
					counter(),
				),

				E("section").C(
					E("h2").A("class", "text-xl font-semibold mb-3").T("Accordion"),
					accordion("What is gomb?", "gomb is a Go library for building HTML programmatically using a fluent API."),
					accordion("What is Alpine.js?", "Alpine.js adds reactive behaviour directly in HTML using x-* attributes."),
					accordion("Do they work together?", "Yes – gomb generates the HTML with Alpine attributes; Alpine handles client-side reactivity."),
				),

				E("section").C(
					E("h2").A("class", "text-xl font-semibold mb-3").T("Tabs"),
					tabPanel(),
				),

				E("section").C(
					E("h2").A("class", "text-xl font-semibold mb-3").T("Modal"),
					modal(),
				),

				E("section").C(
					E("h2").A("class", "text-xl font-semibold mb-3").T("Live Search"),
					liveSearch(languages),
				),
			),
		),
	)
}

// ── server ───────────────────────────────────────────────────────────────────

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		page().Render(w)
	})

	fmt.Println("Alpine.js example running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
