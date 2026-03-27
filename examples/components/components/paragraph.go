package components

import "github.com/ernlel/gomb"

// Paragraph returns a fragment containing one <p> element per text argument.
func Paragraph(texts ...string) gomb.Element {
	paras := gomb.Map(texts, func(text string) gomb.Element {
		return gomb.E("p").T(text)
	})
	return gomb.Fragment(paras...)
}
