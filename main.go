package main

import (
	"fmt"
	"strings"
)

type Element struct {
	tag      string
	attrs    map[string]string
	children []Element
	text     string
}

func (e Element) a(key, value string) Element {
	if e.attrs == nil {
		e.attrs = make(map[string]string)
	}
	e.attrs[key] = value
	return e
}

func (e Element) t(value string) Element {
	e.text = value
	return e
}

func (e Element) c(elements ...Element) Element {
	e.children = append(e.children, elements...)
	return e
}

func (e Element) toString() string {
	var sb strings.Builder
	sb.WriteString("<" + e.tag)
	for k, v := range e.attrs {
		sb.WriteString(fmt.Sprintf(` %s="%s"`, k, v))
	}
	sb.WriteString(">")
	if e.text != "" {
		sb.WriteString(e.text)
	}
	for _, child := range e.children {
		sb.WriteString(child.toString())
	}
	sb.WriteString("</" + e.tag + ">")
	return sb.String()
}

func (e Element) toStringIndented(indent string) string {
	var sb strings.Builder
	sb.WriteString(indent + "<" + e.tag)
	for k, v := range e.attrs {
		sb.WriteString(fmt.Sprintf(` %s="%s"`, k, v))
	}
	sb.WriteString(">\n")
	if e.text != "" {
		sb.WriteString(indent + "  " + e.text + "\n")
	}
	for _, child := range e.children {
		sb.WriteString(child.toStringIndented(indent + "  "))
	}
	sb.WriteString(indent + "</" + e.tag + ">\n")
	return sb.String()
}

func E(tag string) Element {
	return Element{tag: tag}
}

func main() {

	div := E("div").a("id", "my-div").a("style", "background-color:red;").c(
		E("h1").t("Hello"),
		E("ul").c(
			// insert 10 elements
			func() (elements []Element) {
				for i := 0; i < 10; i++ {
					elements = append(elements, E("li").t(fmt.Sprintf("item%d", i)))
				}
				return
			}()...,
		),
		E("ul").c(
			E("li").t("I am true"),
			E("li").t("item1"),
			E("br"),
			E("li").t("item2"),
			func() Element {
				if false {
					return E("li").t("I am true")
				} else if true {
					return E("li").t("I am true second")
				} else {
					return E("li").t("I am true second")
				}
			}(),
		),
	).toStringIndented("")

	fmt.Println(div)

	htmlStr := `
    <div id="my-div" style="background-color:red;">
        <h1>Hello</h1>
        <ul>
            <li>I am true</li>
            <li>item1</li>
            <br>
            <li>item2</li>
            <li>I am true second</li>
        </ul>
    </div>
    `

	goCode, err := generateGoCodeFromHTMLString(htmlStr)
	if err != nil {
		fmt.Println("Error generating Go code:", err)
		return
	}

	fmt.Println(goCode)

	div2 := E("html").c(
		E("head"),
		E("body").c(
			E("div").a("id", "my-div").a("style", "background-color:red;").c(
				E("h1").t("Hello"),
				E("ul").c(
					E("li").t("I am true"),
					E("li").t("item1"),
					E("br"),
					E("li").t("item2"),
					E("li").t("I am true second"),
				),
			),
		),
	).toString()

	fmt.Println(div2)
}
