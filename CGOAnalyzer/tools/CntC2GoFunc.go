package tools

import (
	"database/sql"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"regexp"
)

type C2GoFuncVisitor struct {
	fileName string
	set      *token.FileSet
	db       *sql.DB
	pkgid    int
}

func (v C2GoFuncVisitor) Visit(n ast.Node) ast.Visitor {
	if fc, ok := n.(*ast.FuncDecl); ok {
		funcname := fc.Name.Name
		if fc.Doc != nil {
			var regExportFunc *regexp.Regexp = regexp.MustCompile(`// *export +` + funcname + ` *.*$`)
			for _, comment := range fc.Doc.List {
				if regExportFunc.Match([]byte(comment.Text)) {
					fmt.Println(comment.Text)
				}
			}
		}
	}
	return v
}

func CntC2GoFunc_Pkg(pkgast *ast.Package, pkgid int, set *token.FileSet, db *sql.DB) {
	// 解析文件
	for filename, srcfile := range pkgast.Files {
		msgReport(os.Stdout, showAll, "-- now parsing file:%s\n", filename)
		pkg_set := FilterFile(srcfile, "C")
		if len(pkg_set) == 0 {
			continue
		}

		// 统计函数
		v := C2GoFuncVisitor{pkgid: pkgid, fileName: filename, set: set, db: db}
		ast.Walk(v, srcfile)
	}
}
