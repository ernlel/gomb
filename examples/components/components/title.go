package components

import (
	"fmt"

	"github.com/ernlel/gomb"
)

// Title creates a title element with the specified text and size.
func Title(text string, size int) gomb.Element {
	tag := fmt.Sprintf("h%d", size)
	return gomb.E(tag).T(text)
}
