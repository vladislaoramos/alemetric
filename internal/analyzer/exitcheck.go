// Package exitcheck defines an Analyzer that detects
// fragments of the source code in which the program is terminated immediately.
package exitcheck

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

const Doc = `checks for os.Exit inside the main function

This check allows you to detect fragments of the source code in which the program is terminated immediately.
`

var Analyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  Doc,
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.CallExpr:
				switch ce := x.Fun.(type) {
				case *ast.SelectorExpr:
					if fmt.Sprintf("%s", ce.X) == "os" && ce.Sel.Name == "Exit" {
						pass.Reportf(ce.Pos(), "unwanted use of os.Exit")
					}
				}
			case *ast.FuncDecl:
				if x.Name.Name != "main" {
					return false
				}
			}
			return true
		})
	}

	return nil, nil
}
