# GoMB - Go Markup Builder

`. "github.com/ernlel/gomb"`

- E - Element
- A - Attribute
- T - Text
- C - Children
```
E("div").A("id", "my-div").A("style","background-color:red;").C(
		E("h1").T("Hello"), 
		E("ul").C(
			If(true, 
				E("li").T("I am true")
			),
			If(true, 
				E("li").T("I am true")
			).elseIf(true,
				E("li").T("I am true second")
			).else(
				E("li").T("I am true second")
			),
			E("li").T("item1"), 
			E("li").T("item2")
		)
	)
	```
