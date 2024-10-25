# User Instruction for GOMB Library

## Introduction

GOMB (Go Markup Builder) is a Go library designed to help you create and manipulate HTML elements programmatically. This guide will walk you through the usage of the GOMB library, providing detailed instructions and examples.

## Installation

To install the GOMB library, use the following command:

```sh
go get github.com/yourusername/gomb
```

## Basic Usage

### Creating an Element

To create a new HTML element, use the `E` function:

```go
div := gomb.E("div")
```

### Adding Attributes

To add attributes to an element, use the `A` method:

```go
div = div.A("class", "container")
```

### Adding Text

To add text content to an element, use the `T` method:

```go
div = div.T("Hello, World!")
```

### Adding Child Elements

To add child elements, use the `C` method:

```go
span := gomb.E("span").T("This is a span.")
div = div.C(span)
```

### Generating HTML String

To generate the HTML string representation of an element, use the `ToString` method:

```go
html := div.ToString()
fmt.Println(html)
```

## Example

Here is a complete example that demonstrates how to create a nested HTML structure:

```go
package main

import (
    "fmt"
    "github.com/yourusername/gomb"
)

func main() {
    div := gomb.E("div").
        A("class", "container").
        C(
            gomb.E("h1").T("Welcome to GOMB"),
            gomb.E("p").T("This is a paragraph."),
            gomb.E("ul").C(
                gomb.E("li").T("Item 1"),
                gomb.E("li").T("Item 2"),
                gomb.E("li").T("Item 3"),
            ),
        )

    fmt.Println(div.ToString())
}
```

This will produce the following HTML:

```html
<div class="container">
  <h1>
    Welcome to GOMB
  </h1>
  <p>
    This is a paragraph.
  </p>
  <ul>
    <li>
      Item 1
    </li>
    <li>
      Item 2
    </li>
    <li>
      Item 3
    </li>
  </ul>
</div>
```

## Advanced Usage

### Self-Closing Tags

The GOMB library automatically handles self-closing tags. For example:

```go
img := gomb.E("img").A("src", "image.png").A("alt", "An image")
fmt.Println(img.ToString())
```

This will produce:

```html
<img src="image.png" alt="An image" />
```

### Nested Elements

You can nest elements to any depth:

```go
outerDiv := gomb.E("div").A("class", "outer").C(
    gomb.E("div").A("class", "inner").T("Inner content"),
)
fmt.Println(outerDiv.ToString())
```

This will produce:

```html
<div class="outer">
  <div class="inner">
    Inner content
  </div>
</div>
```

## Conclusion

The GOMB library provides a simple and intuitive way to build and manipulate HTML elements in Go. By following this guide, you should be able to create complex HTML structures programmatically with ease. For more information and advanced usage, refer to the official documentation and examples provided in the repository.
