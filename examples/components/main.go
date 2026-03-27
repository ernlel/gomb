package main

import (
	"fmt"
	"net/http"

	"github.com/ernlel/gomb"
	"github.com/ernlel/gomb/examples/components/components"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	page := gomb.E("html").A("lang", "en").C(
		gomb.E("head").C(
			gomb.E("title").T("Gomb Example"),
		),
		gomb.E("body").C(
			components.Title("Welcome to Gomb", 1),
			components.Paragraph("This is a paragraph.", "Here is another paragraph."),
			components.TextInput("Name:", "Enter your name"),
			components.Button("Submit"),
		),
	)

	fmt.Fprint(w, page.ToString())
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
