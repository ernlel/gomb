package pages

import (
	"github.com/ernlel/gomb"
	"github.com/ernlel/gomb/examples/layout/template/layout"
)

func NotFoundPage() gomb.Element {
	E := gomb.E
	content := E("div").
		A("class", "not-found-content container mx-auto p-4").
		C(
			E("h2").
				A("class", "text-3xl font-bold mb-4").
				T("404 - Page Not Found"),
			E("p").
				A("class", "text-lg mb-4").
				T("Sorry, the page you are looking for does not exist."),
			E("a").
				A("class", "text-blue-500 underline").
				A("href", "/").
				T("Go back to Home Page"),
		)
	return layout.Layout(content)
}
