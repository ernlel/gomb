package pages

import (
	"github.com/ernlel/gomb"
	"github.com/ernlel/gomb/examples/layout/template/layout"
)

func HomePage() *gomb.Element {
	E := gomb.E
	content := E("div").
		A("class", "home-content container mx-auto p-4").
		C(
			E("h2").
				A("class", "text-3xl font-bold mb-4").
				T("Welcome to the Home Page"),
			E("p").
				A("class", "text-lg mb-4").
				T("This is the main content of the home page."),
			E("div").
				A("class", "features mt-8").
				C(
					E("h3").
						A("class", "text-2xl font-semibold mb-2").
						T("Features"),
					E("ul").
						A("class", "list-disc list-inside").
						C(
							E("li").A("class", "mb-2").T("Feature 1: Description of feature 1."),
							E("li").A("class", "mb-2").T("Feature 2: Description of feature 2."),
							E("li").A("class", "mb-2").T("Feature 3: Description of feature 3."),
						),
				),
			E("div").
				A("class", "testimonials mt-8").
				C(
					E("h3").
						A("class", "text-2xl font-semibold mb-2").
						T("Testimonials"),
					E("div").
						A("class", "testimonial-item mb-4").
						C(
							E("p").
								A("class", "text-lg italic").
								T("\"This is a great product!\" - User 1"),
						),
					E("div").
						A("class", "testimonial-item mb-4").
						C(
							E("p").
								A("class", "text-lg italic").
								T("\"I highly recommend this.\" - User 2"),
						),
				),
			E("div").
				A("class", "contact mt-8").
				C(
					E("h3").
						A("class", "text-2xl font-semibold mb-2").
						T("Contact Us"),
					E("p").
						A("class", "text-lg").
						T("For more information, please contact us at info@example.com."),
				),
		)
	return layout.Layout(content)
}
