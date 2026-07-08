// Command gomblint checks for ignored return values from gomb builder methods.
//
// gomb uses immutable value types: every A(), T(), C(), and helper method returns
// a new Element. Discarding the return value is always a bug:
//
//	body.C(form)          // BUG — copy is discarded, form is never added
//	body = body.C(form)   // CORRECT — reassign
//
// Usage:
//
//	go run ./cmd/gomblint ./...
//
// Run in CI alongside your regular tests.
package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/ast/inspector"
)

var gombMethods = map[string]bool{
	"A": true, "T": true, "C": true,
	"Attr": true, "Text": true, "Children": true,
	"Attrs": true, "As": true,
	"Data": true, "Style": true,
	"With": true, "Render": true,
}

var Analyzer = &analysis.Analyzer{
	Name:     "gomblint",
	Doc:      "reports ignored return values from gomb builder methods",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filter to expression statements only.
	nodeFilter := []ast.Node{(*ast.ExprStmt)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		stmt := n.(*ast.ExprStmt)

		call, ok := stmt.X.(*ast.CallExpr)
		if !ok {
			return
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		methodName := sel.Sel.Name
		if !gombMethods[methodName] {
			return
		}

		// Check receiver is gomb.Element
		recvType := pass.TypesInfo.TypeOf(sel.X)
		if recvType == nil {
			return
		}
		if !isGombElement(recvType) {
			return
		}

		pass.Reportf(call.Pos(),
			"return value of %s() is discarded — gomb methods return a new Element; reassign or chain: el = el.%s(...)",
			methodName, methodName)
	})

	return nil, nil
}

func isGombElement(t interface{ String() string }) bool {
	s := t.String()
	return s == "github.com/ernlel/gomb.Element" || s == "*github.com/ernlel/gomb.Element"
}

func main() {
	singlechecker.Main(Analyzer)
}
