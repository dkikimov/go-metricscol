package osexitchecker

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OsExitChecker checks if os.Exit is being used in main function in main package.
var OsExitChecker = &analysis.Analyzer{
	Name: "osexitchecker",
	Doc:  "check for os.Exit usage in main function in main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}

		for _, decl := range file.Decls {
			mainDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if mainDecl.Name.Name != "main" {
				continue
			}

			ast.Inspect(mainDecl, func(node ast.Node) bool {
				if c, ok := node.(*ast.CallExpr); ok {
					if s, ok := c.Fun.(*ast.SelectorExpr); ok {
						if xIdent, ok := s.X.(*ast.Ident); ok {
							if xIdent.Name == "os" && s.Sel.Name == "Exit" {
								pass.Reportf(s.Sel.End(), "using os.Exit in main file of main package")
							}
						}
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
