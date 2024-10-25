package components

import "github.com/ernlel/gomb"

func Paragraph(texts ...string) gomb.Element {
	container := gomb.E("")
	for _, text := range texts {
		container.Children = append(container.Children, gomb.E("p").T(text))
	}
	return container
}
