package main

import (
	"fmt"
	"net/http"

	"github.com/ernlel/gomb/examples/layout/template/pages"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			w.Write([]byte(pages.HomePage().ToString()))
		case "/about":
			w.Write([]byte(pages.AboutPage().ToString()))
		case "/contact":
			w.Write([]byte(pages.ContactPage().ToString()))
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(pages.NotFoundPage().ToString()))
		}
	})

	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
