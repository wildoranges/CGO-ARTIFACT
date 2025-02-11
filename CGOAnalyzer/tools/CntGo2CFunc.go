package tools

import (
	"database/sql"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type FuncItem struct {
	label string
}

const GoCTypeTrans int = 0
const StdCLibFunc int = 1
const FuncUnknown int = 2
const StdCType int = 3

// var pkg_id int = 0

var FuncKind = map[int]string{
	GoCTypeTrans: "GoCTypeTrans",
	StdCLibFunc:  "StdCLibFunc",
	FuncUnknown:  "Unknown",
	StdCType:     "StdCType",
}

func (i FuncItem) KindText(name string) string {
	if _, ok := cStdType[name]; ok {
		return FuncKind[StdCType]
	} else if _, ok := TypeTrans[name]; ok {
		return FuncKind[GoCTypeTrans]
	}
	return FuncKind[FuncUnknown]
}

func (i FuncItem) Label() string {
	return i.label
}

type void struct {
}

var member void

var TypeTrans = map[string]void{
	"GoString": member,
	"CString":  member,
	"GoBytes":  member,
}

var cStdType = map[string]void{ //only part of the std c ctypes and macros
	"schar":         member,
	"uchar":         member,
	"ushort":        member,
	"uint":          member,
	"ulong":         member,
	"longlong":      member,
	"ulonglong":     member,
	"complexfloat":  member,
	"complexdouble": member,
	"float":         member,
	"double":        member,
	"short":         member,
	"int":           member,
	"char":          member,
	"size_t":        member,
	"intptr_t":      member,
	"uintptr_t":     member,
	"u_char":        member,
	"u_short":       member,
	"u_int":         member,
	"u_long":        member,
	"quad_t":        member,
	"u_quad_t":      member,
	"uint_t":        member,
	"int8_t":        member,
	"int16_t":       member,
	"int32_t":       member,
	"int64_t":       member,
	"uint8_t":       member,
	"uint16_t":      member,
	"uint32_t":      member,
	"uint64_t":      member,
	"bool":          member,
	"gint":          member,
	"Uint8":         member,
	"Sint16":        member,
	"Uint32":        member,
	"guint":         member,
	"gdouble":       member,
	"gpointer":      member,
	"long":          member,
	"gfloat":        member,
	"Py_ssize_t":    member,
	"UINT":          member,
}

//GoCFuncVisitor contains the number of each function while parsing a file/repo
type GoCFuncVisitor struct {
	fileName string
	set      *token.FileSet
	db       *sql.DB
	fileID   int
}

var regStd *regexp.Regexp = regexp.MustCompile(`^.*#include *(<.+\.h>).*$`)
var regUser *regexp.Regexp = regexp.MustCompile(`^.*#include *(".+\.h").*$`)
var regLib *regexp.Regexp = regexp.MustCompile(`^.*#cgo LDFLAGS:.*(-l.+).*$`)

func isCStdType(name string) bool {
	_, ok := cStdType[name]
	return ok
}

//Visit counts the number of each function while walking traverse an AST
func (v GoCFuncVisitor) Visit(n ast.Node) ast.Visitor {
	if call, ok := n.(*ast.CallExpr); ok {
		if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
			funcName := fun.Sel.Name
			if pkgname, ok := fun.X.(*ast.Ident); ok {
				msgReport(os.Stdout, showAll, "now checking func:%s\n", funcName)
				if pkgname.Name == "C" {
					if !isCStdType(funcName) {
						// C function
						msgReport(os.Stdout, showAll, "Go2CFunc:%s\n", funcName)

						// get position
						line := v.set.Position(fun.Pos()).Line
						column := v.set.Position(fun.Pos()).Column
						// insert to cgo function table
						funcdata := []string{strconv.Itoa(line), strconv.Itoa(v.fileID), funcName, "", strconv.Itoa(column)}
						Insert2Tabel(v.db, "cgo_function", funcdata)
					}
				}
			}
		}
	} else if fc, ok := n.(*ast.FuncDecl); ok {
		funcname := fc.Name.Name
		if fc.Doc != nil {
			var regExportFunc *regexp.Regexp = regexp.MustCompile(`// *export +` + funcname + ` *$`)
			for _, comment := range fc.Doc.List {
				if regExportFunc.Match([]byte(comment.Text)) {
					line := v.set.Position(fc.Pos()).Line
					column := v.set.Position(fc.Pos()).Column
					exportFuncData := []string{strconv.Itoa(line), strconv.Itoa(v.fileID), funcname, strconv.Itoa(column)}
					Insert2Tabel(v.db, "export_function", exportFuncData)
				}
			}
		}
	}
	return v
}

// 对于每条注释，记录所使用的C头文件和链接库
func RecordHeaderLib(c *ast.CommentGroup, set *token.FileSet, fileID int, db *sql.DB, filename string) {
	// parse header and lib
	cmt := c.Text()
	lines := strings.Split(cmt, "\n")
	for _, line := range lines {
		// 统计lib信息
		libstr := regLib.FindStringSubmatch(line)
		if len(libstr) > 1 {
			line := set.Position(c.Pos()).Line
			collomn := set.Position(c.Pos()).Column
			libs := strings.Split(libstr[1], "-l")
			for _, lib := range libs {
				libData := []string{strconv.Itoa(line), strconv.Itoa(fileID), lib, "", strconv.Itoa(collomn)}
				err := Insert2Tabel(db, "cgo_lib", libData)
				if err != nil {
					panic(err)
				}
			}
		}

		// 统计header信息
		var headers []string
		headers = regStd.FindStringSubmatch(line)
		if len(headers) <= 1 {
			headers = regUser.FindStringSubmatch(line)
		}
		if len(headers) > 1 {
			header := headers[1]
			msgReport(os.Stdout, showAll, "found include header %s in file %s\n", header, filepath.Base(filename))
			line := set.Position(c.Pos()).Line
			collomn := set.Position(c.Pos()).Column
			// insert to header table
			includeData := []string{strconv.Itoa(line), strconv.Itoa(fileID), header, "", strconv.Itoa(collomn)}
			errs := Insert2Tabel(db, "cgo_header", includeData)
			if errs != nil {
				fmt.Println(includeData)
				panic(errs)
			}
		}
	}
}

// 记录Go调用C的函数
func RecordGo2CFunc(filename string, set *token.FileSet, db *sql.DB, fileID int, srcfile *ast.File) {
	m := GoCFuncVisitor{fileName: filename, set: set, db: db, fileID: fileID}
	msgReport(os.Stdout, showAll, "now parsing file:%s\n", filename)
	ast.Walk(m, srcfile)
}

func CntGo2CFunc_File(srcfile *ast.File, fileID int, set *token.FileSet, db *sql.DB, filename string) {
	// 当前文件是否含有import"C"
	pkg_set := FilterFile(srcfile, "C")
	if len(pkg_set) == 0 {
		return
	}

	// parse header and lib
	for _, c := range srcfile.Comments {
		RecordHeaderLib(c, set, fileID, db, filename)
	}

	// parse function
	RecordGo2CFunc(filename, set, db, fileID, srcfile)
}

// func CntGo2CFunc_Pkg(pkgast *ast.Package, pkgid int, set *token.FileSet, db *sql.DB) {
// 	for filename, srcfile := range pkgast.Files {
// 		pkg_set := FilterFile(srcfile, "C")
// 		if len(pkg_set) == 0 {
// 			continue
// 		}

// 		// parse header and lib
// 		for _, c := range srcfile.Comments {
// 			cmt := c.Text()
// 			lines := strings.Split(cmt, "\n")
// 			for _, line := range lines {
// 				// 统计lib信息
// 				libstr := regLib.FindStringSubmatch(line)
// 				if len(libstr) > 1 {
// 					line := set.Position(c.Pos()).Line
// 					collomn := set.Position(c.Pos()).Column
// 					libs := strings.Split(libstr[1], "-l")
// 					for _, lib := range libs {
// 						libData := []string{strconv.Itoa(line), strconv.Itoa(pkgid), lib, "", strconv.Itoa(collomn)}
// 						err := Insert2Tabel(db, "cgo_lib", libData)
// 						if err != nil {
// 							panic(err)
// 						}
// 					}
// 				}

// 				// 统计header信息
// 				var headers []string
// 				headers = regStd.FindStringSubmatch(line)
// 				if len(headers) <= 1 {
// 					headers = regUser.FindStringSubmatch(line)
// 				}
// 				if len(headers) > 1 {
// 					header := headers[1]
// 					msgReport(os.Stdout, showAll, "found include header %s in file %s\n", header, filepath.Base(filename))
// 					line := set.Position(c.Pos()).Line
// 					collomn := set.Position(c.Pos()).Column
// 					// insert to header table
// 					includeData := []string{strconv.Itoa(line), strconv.Itoa(pkgid), header, "", strconv.Itoa(collomn)}
// 					errs := Insert2Tabel(db, "cgo_header", includeData)
// 					if errs != nil {
// 						fmt.Println(includeData)
// 						panic(errs)
// 					}
// 				}

// 			}
// 		}

// 		// parse function
// 		curfile := make(map[string]int)
// 		m := GoCFuncVisitor{curFile: curfile, fileName: filename, set: set, db: db, pkgid: pkgid}
// 		msgReport(os.Stdout, showAll, "now parsing file:%s\n", filename)
// 		ast.Walk(m, srcfile)
// 	}
// }

// func CntGo2CInfo_Repo(repo string, db *sql.DB) {

// 	msgReport(os.Stdout, showAll, "now parsing repo:%s\n", repo)

// 	// 获取当前repo对应的id
// 	repoName := path.Base(repo)
// 	id, _ := SelectRepoID(db, repoName)

// 	alldirs, err := GetAllDirs(repo)
// 	if err != nil {
// 		msgReport(os.Stderr, true, err.Error()+"\n")
// 		return
// 	}

// 	for _, repoDir := range alldirs {
// 		set := token.NewFileSet()
// 		f, err := parser.ParseDir(set, repoDir, nil, parser.ParseComments)
// 		if err != nil {
// 			msgReport(os.Stderr, true, err.Error()+"\n")
// 		}
// 		for pkg, pkgast := range f {
// 			// insert to pakcage table
// 			pkg_data := []string{strconv.Itoa(pkg_id), pkg, strconv.Itoa(id), repoDir}
// 			Insert2Tabel(db, "package", pkg_data)
// 			pkg_id += 1

// 			msgReport(os.Stdout, showAll, "now parsing pkg:%s\n", pkg)
// 			for filename, srcfile := range pkgast.Files {
// 				// parse import
// 				for _, c := range srcfile.Comments {
// 					cmt := c.Text()
// 					lines := strings.Split(cmt, "\n")
// 					for _, line := range lines {
// 						// 统计lib信息
// 						libstr := regLib.FindStringSubmatch(line)
// 						if len(libstr) > 1 {
// 							line := set.Position(c.Pos()).Line
// 							collomn := set.Position(c.Pos()).Column
// 							libs := strings.Split(libstr[1], "-l")
// 							for _, lib := range libs {
// 								libData := []string{filename, strconv.Itoa(line), strconv.Itoa(pkg_id), lib, "", strconv.Itoa(collomn)}
// 								err := Insert2Tabel(db, "cgo_lib", libData)
// 								if err != nil {
// 									panic(err)
// 								}
// 							}
// 						}

// 						// 统计header信息
// 						var headers []string
// 						headers = regStd.FindStringSubmatch(line)
// 						if len(headers) <= 1 {
// 							headers = regUser.FindStringSubmatch(line)
// 						}
// 						if len(headers) > 1 {
// 							header := headers[1]
// 							msgReport(os.Stdout, showAll, "found include header %s in file %s\n", header, filepath.Base(filename))
// 							line := set.Position(c.Pos()).Line
// 							collomn := set.Position(c.Pos()).Column
// 							// insert to header table
// 							includeData := []string{filename, strconv.Itoa(line), strconv.Itoa(pkg_id), header, "", strconv.Itoa(collomn)}
// 							errs := Insert2Tabel(db, "cgo_header", includeData)
// 							if errs != nil {
// 								fmt.Println(includeData)
// 								panic(errs)
// 							}
// 						}

// 					}
// 				}

// 				curfile := make(map[string]int)
// 				m := GoCFuncVisitor{curFile: curfile, fileName: filename, set: set, db: db}
// 				msgReport(os.Stdout, showAll, "now parsing file:%s\n", filename)
// 				ast.Walk(m, srcfile)
// 			}
// 		}
// 	}
// }

// //CntGo2CFunc counts the go-c function in each repo,return two maps.
// //the first map contains the overall result.its key is C function name and its value is the number of the function
// //the second map contains the result in each repo,the key is repo's name and the value is a map containing
// //the repo's local result
// func CntGo2CFunc(repos []string) {
// 	db, _ := ConnectSQL()

// 	for _, repo := range repos {
// 		CntGo2CInfo_Repo(repo, db)
// 	}

// 	db.Close()
// }
