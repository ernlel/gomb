package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ernlel/gomb/pkg/markup_to_gomb"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main <filename>")
		return
	}

	filename := os.Args[1]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	gombCode, err := markup_to_gomb.GenerateGombFromMarkup(string(content))
	if err != nil {
		fmt.Printf("Error generating Go code: %v\n", err)
		return
	}

	goCode := `
package main

import (
	"fmt"
	"os"

	"github.com/ernlel/gomb"
)

var E = gomb.E

var template =  ` + gombCode + `

var html = template.ToString()

func main() {
	// save to file
	f, err := os.Create("generated_index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	f.WriteString(html)
	fmt.Println("HTML file generated successfully")
}
	`

	newFilename := changeExtensionToGo(filename)
	err = os.WriteFile("output/"+newFilename, []byte(goCode), 0o644)
	if err != nil {
		fmt.Printf("Error writing Go code to file: %v\n", err)
		return
	}

	fmt.Printf("Go code successfully written to %s\n", newFilename)
}

func changeExtensionToGo(filename string) string {
	ext := filepath.Ext(filename)
	return filename[:len(filename)-len(ext)] + ".go"
}
