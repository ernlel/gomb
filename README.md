# gomb
Go Markup Builder

Node("div").attr("id", "my-div").attr("style","background-color:red;").children(
		Node("h1").text("Hello"), 
		Node("ul").children(
			If(true, 
				Node("li").text("I am true")
			),
			If(true, 
				Node("li").text("I am true")
			).elseIf(true,
				Node("li").text("I am true second")
			).else(
				Node("li").text("I am true second")
			),

			

			Node("li").text("item1"), 
			Node("li").text("item2")
		)
	)
