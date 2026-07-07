// Package html provides named constructor functions for every standard HTML
// element. Each function returns a gomb.Element, giving you IDE autocomplete,
// compile-time safety, and GoDoc on every HTML tag.
//
// Import as a dot-import for a natural DSL:
//
//	import . "github.com/ernlel/gomb"
//	import . "github.com/ernlel/gomb/html"
//
//	page := Html(
//	    Head(Meta(Attr{Key: "charset", Value: "UTF-8"})),
//	    Body(H1(Txt("Hello")), P(Txt("Welcome"))),
//	).Attr("lang", "en")
//
// Or with a qualified import:
//
//	import "github.com/ernlel/gomb"
//	import "github.com/ernlel/gomb/html"
//
//	html.Div(html.Attr{Key: "class", Value: "box"}, html.P(html.Txt("hi")))
package html

import "github.com/ernlel/gomb"

// buildElement is the shared implementation for all named constructors.
// It processes the variadic args: Attr values set attributes, Element values
// become children.
func buildElement(tag string, args []interface{}) gomb.Element {
	el := gomb.E(tag)
	for _, arg := range args {
		switch v := arg.(type) {
		case gomb.Attr:
			el = el.A(v.Key, v.Value)
		case gomb.Element:
			el = el.C(v)
		}
	}
	return el
}
