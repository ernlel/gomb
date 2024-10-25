package pages

import (
	"github.com/ernlel/gomb"
	"github.com/ernlel/gomb/examples/layout/template/layout"
)

func ContactPage() gomb.Element {
	E := gomb.E
	content := E("div").
		A("class", "contact-page container mx-auto p-4").
		C(
			E("h2").
				A("class", "text-3xl font-bold mb-4").
				T("Contact Us"),
			E("form").
				A("action", "/submit-contact").
				A("method", "post").
				A("class", "space-y-4").
				C(
					E("div").
						A("class", "form-group").
						C(
							E("label").
								A("for", "name").
								A("class", "block text-sm font-medium text-gray-700").
								T("Name:"),
							E("input").
								A("type", "text").
								A("id", "name").
								A("name", "name").
								A("class", "mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"),
						),
					E("div").
						A("class", "form-group").
						C(
							E("label").
								A("for", "email").
								A("class", "block text-sm font-medium text-gray-700").
								T("Email:"),
							E("input").
								A("type", "email").
								A("id", "email").
								A("name", "email").
								A("class", "mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"),
						),
					E("div").
						A("class", "form-group").
						C(
							E("label").
								A("for", "message").
								A("class", "block text-sm font-medium text-gray-700").
								T("Message:"),
							E("textarea").
								A("id", "message").
								A("name", "message").
								A("rows", "4").
								A("class", "mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"),
						),
					E("div").
						A("class", "form-group").
						C(
							E("button").
								A("type", "submit").
								A("class", "inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500").
								T("Submit"),
						),
				),
		)

	return layout.Layout(content)
}
