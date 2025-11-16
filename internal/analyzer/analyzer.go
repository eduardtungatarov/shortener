package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// HardStopAnalyzer анализатор жестких остановок.
var HardStopAnalyzer = &analysis.Analyzer{
	Name: "hardStopAnalyzer",
	Doc:  "check for hard stop instruction",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	pkgName := pass.Pkg.Name()
	for _, file := range pass.Files {
		inMain := false
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if pkgName == "main" && x.Name.Name == "main" {
					inMain = true
				} else {
					inMain = false
				}
			case *ast.Ident:
				if x.Name == "panic" {
					pass.Reportf(x.Pos(), "usage of panic")
				}
			case *ast.SelectorExpr:
				funcName := x.Sel.Name
				if ident, ok := x.X.(*ast.Ident); ok {
					pkgName := ident.Name
					if !inMain && pkgName == "log" && funcName == "Fatal" {
						pass.Reportf(x.Pos(), "usage of log.Fatal()")
					}
					if !inMain && pkgName == "os" && funcName == "Exit" {
						pass.Reportf(x.Pos(), "usage of os.Exit()")
					}
				}
			}
			return true
		})
	}

	return nil, nil
}
