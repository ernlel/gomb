package components

import (
	"github.com/ernlel/gomb"
)

func Button(label string) *gomb.Element {
	button := gomb.E("button").A("type", "submit").T(label)
	return button
}
