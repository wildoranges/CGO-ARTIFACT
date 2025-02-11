package tools

import (
	"fmt"
	"go/ast"
)

func CntImport(file *ast.File) {
	for _, imp := range file.Imports {
		fmt.Println(imp.Path.Value)
	}
}
