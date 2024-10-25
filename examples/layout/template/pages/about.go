package pages

import (
	"github.com/ernlel/gomb"
	"github.com/ernlel/gomb/examples/layout/template/layout"
)

func AboutPage() gomb.Element {
	E := gomb.E
	content := E("div").
		A("class", "about-page container mx-auto p-4").
		C(
			E("h2").
				A("class", "text-3xl font-bold mb-4").
				T("About Us"),
			E("p").
				A("class", "text-lg").
				T("This is the about page content."),
		)
	return layout.Layout(content)
}
