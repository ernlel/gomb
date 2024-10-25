package components

import "github.com/ernlel/gomb"

func TextInput(label, placeholder string) gomb.Element {
	div := gomb.E("div")
	if label != "" {
		div = div.C(gomb.E("label").T(label))
	}
	div = div.C(gomb.E("input").A("type", "text").A("placeholder", placeholder))
	return div
}
