package layout

import "github.com/ernlel/gomb"

func Footer() *gomb.Element {
	E := gomb.E
	return E("footer").
		A("class", "bg-gray-800 text-white p-4").
		C(
			E("div").
				A("class", "footer-content text-center").
				T("This is the footer content."),
			E("div").
				A("class", "footer-links mt-4 flex justify-center").
				C(
					E("a").
						A("href", "/privacy").
						A("class", "text-white hover:text-gray-400 mx-2").
						T("Privacy Policy"),
					E("a").
						A("href", "/terms").
						A("class", "text-white hover:text-gray-400 mx-2").
						T("Terms of Service"),
				),
		)
}
