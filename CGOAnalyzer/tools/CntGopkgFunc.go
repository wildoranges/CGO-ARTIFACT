package tools

import (
	"database/sql"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strconv"
	"strings"
)

type GopkgFuncVisitor struct {
	pkg_set StringSet
	set     *token.FileSet
	db      *sql.DB
	name    string
	fileID  int
}

func (v GopkgFuncVisitor) Visit(n ast.Node) ast.Visitor {
	if call, ok := n.(*ast.CallExpr); ok {
		if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
			funcName := fun.Sel.Name
			if pkgname, ok := fun.X.(*ast.Ident); ok {
				// msgReport(os.Stdout, showAll, "now checking func:%s\n", funcName)
				if _, ok := v.pkg_set[pkgname.Name]; ok {
					// get position
					line := v.set.Position(fun.Pos()).Line
					column := v.set.Position(fun.Pos()).Column
					msgReport(os.Stdout, showAll, "--- %s Function found:%s at %d:%d\n", pkgname.Name, funcName, line, column)
					// insert to function table
					funcdata := []string{strconv.Itoa(line), strconv.Itoa(v.fileID), funcName, pkgname.Name, strconv.Itoa(column)}
					table := v.name + "_invocation"
					err := Insert2Tabel(v.db, table, funcdata)
					if err != nil {
						fmt.Println(strconv.Itoa(v.fileID))
						panic(err)
					}
				}
			}
		}
	}
	return v
}

// FilterFile: 获取所有pkg，如：
// 对于crypto包，我们需要捕获import "crypto"和类似与import "crypto/aes"的包引用
// 并且将它们记录在pkg_set中，此时需要记录的是crypto和aes
func FilterFile(file *ast.File, pkgname string) StringSet {
	pkg_set := make(StringSet)
	for _, s := range file.Imports {
		import_name := strings.Trim(s.Path.Value, "\"")
		if import_name == pkgname {
			pkg_set.Insert(pkgname)
		} else if strings.HasPrefix(import_name, pkgname+"/") {
			pkgs := strings.Split(import_name, "/")
			pkg_set.Insert(pkgs[len(pkgs)-1])
		}
	}
	return pkg_set
}

func CntGopkgFunc_File(pkgname string, fileAst *ast.File, fileID int, set *token.FileSet, db *sql.DB) {
	// for _, pkgname := range pkgNames {
	pkg_set := FilterFile(fileAst, pkgname)
	if len(pkg_set) == 0 {
		return
	}
	// 统计函数
	v := GopkgFuncVisitor{pkg_set: pkg_set, set: set, db: db, name: pkgname, fileID: fileID}
	ast.Walk(v, fileAst)
	// }
}

// func CntGopkgFunc_Pkg(pkgname string, pkgast *ast.Package, pkgid int, set *token.FileSet, db *sql.DB) {
// 	// 解析文件
// 	for filename, srcfile := range pkgast.Files {
// 		msgReport(os.Stdout, showAll, "-- now parsing file:%s\n", filename)
// 		pkg_set := FilterFile(srcfile, pkgname)
// 		if len(pkg_set) == 0 {
// 			continue
// 		}

// 		// 统计函数
// 		v := GopkgFuncVisitor{pkg_set: pkg_set, fileName: filename, set: set, db: db, name: pkgname, pkgid: pkgid}
// 		ast.Walk(v, srcfile)
// 	}
// }

// func CntGopkgFunc(repos []string, pkgname string) {
// 	db, _ := ConnectSQL()
// 	for _, repo := range repos {
// 		CntGopkgFunc_Repo(repo, pkgname, db)
// 	}
// 	db.Close()
// }

// func CntGopkgFunc_Repo(repo string, pkgname string, db *sql.DB) {
// 	msgReport(os.Stdout, showAll, "now parsing repo:%s\n", repo)
// 	// 获取当前repo对应的id和类型，如果已经有该类型则不再重复赋值
// 	// repoName := path.Base(repo)
// 	// id, repo_type := SelectRepoID(db, repoName)
// 	// repo_types := strings.Split(repo_type, ";")
// 	// hasPkgType := false
// 	// for _, t := range repo_types {
// 	// 	if t == pkgname {
// 	// 		hasPkgType = true
// 	// 	}
// 	// }
// 	// if !hasPkgType {
// 	// 	err := UpdateRepoType(db, pkgname, repoName)
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// }

// 	alldirs, err := GetAllDirs(repo)
// 	if err != nil {
// 		msgReport(os.Stderr, true, err.Error()+"\n")
// 		return
// 	}

// 	// 解析每个目录
// 	for _, repoDir := range alldirs {
// 		set := token.NewFileSet()
// 		f, err := parser.ParseDir(set, repoDir, nil, 0)
// 		if err != nil {
// 			msgReport(os.Stderr, true, err.Error()+"\n")
// 		}

// 		// 解析package
// 		for pkg, pkgast := range f {
// 			CntGopkgFunc_Pkg(pkgname, pkg, pkgast, pkg_id, set, db)
// 		}
// 	}
// }
