package layout

import (
	"github.com/ernlel/gomb"
)

func Layout(content gomb.Element) gomb.Element {
	E := gomb.E
	page := E("html").C(
		E("head").C(
			E("title").T("Layout Example"),
			E("link").A("href", "https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css").A("rel", "stylesheet"),
		),
		E("body").A("class", "bg-gray-100 text-gray-900").C(
			Header(),
			content,
			Footer(),
		),
	)
	return page
}
