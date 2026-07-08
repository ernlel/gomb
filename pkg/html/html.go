package html

import "github.com/ernlel/gomb"

func buildElement(tag string, args []interface{}) *gomb.Element {
	el := gomb.E(tag)
	for _, arg := range args {
		switch v := arg.(type) {
		case gomb.Attr:
			el.A(v.Key, v.Value)
		case *gomb.Element:
			el.C(v)
		}
	}
	return el
}
