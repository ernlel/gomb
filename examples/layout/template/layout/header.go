package layout

import "github.com/ernlel/gomb"

func Header() *gomb.Element {
	E := gomb.E
	return E("header").
		A("class", "bg-blue-600 text-white p-4").
		C(
			E("div").
				A("class", "flex items-center justify-between").
				C(
					E("div").
						A("class", "text-2xl font-bold").
						T("Logo"),
					E("h1").
						A("class", "text-xl").
						T("Title"),
					E("nav").
						A("class", "flex space-x-4").
						C(
							E("a").
								A("href", "/").
								A("class", "text-white hover:text-gray-200").
								T("Home"),
							E("a").
								A("href", "/about").
								A("class", "text-white hover:text-gray-200").
								T("About"),
							E("a").
								A("href", "/contact").
								A("class", "text-white hover:text-gray-200").
								T("Contact"),
						),
				),
		)
}
