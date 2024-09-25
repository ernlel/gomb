package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func generateGoCodeFromHTMLNode(n *html.Node, indentLevel int) string {
	var sb strings.Builder
	indent := strings.Repeat("    ", indentLevel)

	if n.Type == html.TextNode {
		trimmedText := strings.TrimSpace(n.Data)
		if trimmedText != "" {
			sb.WriteString(fmt.Sprintf("%sE(\"\").t(\"%s\")", indent, trimmedText))
		}
	} else if n.Type == html.ElementNode {
		sb.WriteString(fmt.Sprintf("%sE(\"%s\")", indent, n.Data))
		for _, attr := range n.Attr {
			sb.WriteString(fmt.Sprintf(".a(\"%s\", \"%s\")", attr.Key, attr.Val))
		}

		// Check if the element has only one text child
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode && n.FirstChild.NextSibling == nil {
			trimmedText := strings.TrimSpace(n.FirstChild.Data)
			if trimmedText != "" {
				sb.WriteString(fmt.Sprintf(".t(\"%s\")", trimmedText))
			}
		} else if n.FirstChild != nil {
			sb.WriteString(".c(\n")
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				childCode := generateGoCodeFromHTMLNode(c, indentLevel+1)
				if childCode != "" {
					sb.WriteString(childCode + ",\n")
				}
			}
			sb.WriteString(fmt.Sprintf("%s)", indent))
		}
	}

	return sb.String()
}

func generateGoCodeFromHTMLString(htmlStr string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", err
	}

	// Find the root element (skip the document node)
	var root *html.Node
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			root = c
			break
		}
	}

	if root == nil {
		return "", fmt.Errorf("no root element found")
	}

	return generateGoCodeFromHTMLNode(root, 0), nil
}
