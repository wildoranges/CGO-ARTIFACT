package tools

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func isExported(name string) bool {
	if name[0] >= 'A' && name[0] <= 'Z' {
		return true
	} else {
		return false
	}
}

// 该函数返回一个文件夹中出现的Go函数（包括函数名称、参数、返回值）、全局变量、类型
func CntGoFunc(dir string) {
	// 定义表头
	func_header := []string{"package", "file", "line", "column", "name", "param amount", "params", "return-value amount", "return-values", "function declaration"}
	gvar_header := []string{"package", "file", "line", "column", "name", "type"}
	type_header := []string{"package", "file", "line", "column", "name"}
	if _, err := os.Stat("../data/func.csv"); os.IsNotExist(err) {
		WriteCSVData("../data/func.csv", func_header)
	}
	if _, err := os.Stat("../data/gvar.csv"); os.IsNotExist(err) {
		WriteCSVData("../data/gvar.csv", gvar_header)
	}
	if _, err := os.Stat("../data/type.csv"); os.IsNotExist(err) {
		WriteCSVData("../data/type.csv", type_header)
	}

	// 定义数据表
	var func_dataset, type_dataset, gvar_dataset [][]string

	set := token.NewFileSet()
	f, err := parser.ParseDir(set, dir, nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	for pkg, pkgAST := range f {
		if !strings.HasSuffix(pkg, "_test") {
			func_dataset = make([][]string, 0)
			type_dataset = make([][]string, 0)
			gvar_dataset = make([][]string, 0)
			for filename, fileAST := range pkgAST.Files {
				if strings.HasSuffix(filename, "_test.go") {
					continue
				}
				for _, decl := range fileAST.Decls {
					// 统计函数信息
					if function, ok := decl.(*ast.FuncDecl); ok {
						if isExported(function.Name.Name) && function.Recv == nil {
							line := set.Position(function.Pos()).Line
							column := set.Position(function.Pos()).Column
							relpath, _ := filepath.Rel("/home/dby/go-c/new-go/go", filename)
							relpath = "https://github.com/golang/go/blob/go1.16.3/" + relpath
							relpath += "#L" + strconv.Itoa(line)
							// 初始化func_data
							func_data := []string{pkg, relpath, strconv.Itoa(line), strconv.Itoa(column), function.Name.Name}

							// 统计参数信息
							var param_count int = 0
							var param_str string = ""
							if function.Type.Params != nil {
								for _, param := range function.Type.Params.List {
									param_count += len(param.Names)
									for _, param_name := range param.Names {
										param_str += param_name.Name
										param_str += "; "
										// param_count += 1
									}
									ident, ok := param.Type.(*ast.Ident)
									if ok {
										param_str += fmt.Sprintf(" Type: %v.    ", ident.Name)
									} else {
										param_str += fmt.Sprintf(" Type: %v.    ", param.Type)
									}
								}
							}
							func_data = append(func_data, strconv.Itoa(param_count))
							func_data = append(func_data, param_str)

							// 统计返回值信息
							var result_count int = 0
							var result_str string = ""
							if function.Type.Results != nil {
								for _, result := range function.Type.Results.List {
									result_count += len(result.Names)
									for _, result_name := range result.Names {
										result_str += result_name.Name
										result_str += "; "
										result_count += 1
									}
									if len(result.Names) == 0 {
										result_count += 1
									}
									ident, ok := result.Type.(*ast.Ident)
									if ok {
										result_str += fmt.Sprintf(" Type: %v.    ", ident.Name)
									} else {
										result_str += fmt.Sprintf(" Type: %v.    ", result.Type)
									}
								}
							}
							func_data = append(func_data, strconv.Itoa(result_count))
							func_data = append(func_data, result_str)

							// 获取函数定义
							f, err := os.Open(filename)
							if err != nil {
								log.Fatal(err)
							}
							defer f.Close()
							reader := bufio.NewReader(f)
							count := 0
							for {
								count += 1
								line, err := reader.ReadString('\n')
								if err != nil {
									log.Fatal(err)
								}
								if count == set.Position(function.Pos()).Line {
									line = line[0 : len(line)-2]
									func_data = append(func_data, line)
									break
								}
							}

							func_dataset = append(func_dataset, func_data)
						}
					}

					if gendecl, ok := decl.(*ast.GenDecl); ok {
						for _, spec := range gendecl.Specs {
							// 统计类型信息
							if typespec, ok := spec.(*ast.TypeSpec); ok {
								if isExported(typespec.Name.Name) {
									line := set.Position(typespec.Pos()).Line
									column := set.Position(typespec.Pos()).Column
									relpath, _ := filepath.Rel("/home/dby/go-c/new-go/go", filename)
									relpath = "https://github.com/golang/go/blob/go1.16.3/" + relpath
									relpath += "#L" + strconv.Itoa(line)
									type_data := []string{pkg, relpath, strconv.Itoa(line), strconv.Itoa(column), typespec.Name.Name}
									type_dataset = append(type_dataset, type_data)
								}
							}
							// 统计变量信息
							if valuespec, ok := spec.(*ast.ValueSpec); ok {
								for _, name := range valuespec.Names {
									if isExported(name.Name) {
										line := set.Position(valuespec.Pos()).Line
										column := set.Position(valuespec.Pos()).Column
										relpath, _ := filepath.Rel("/home/dby/go-c/new-go/go", filename)
										relpath = "https://github.com/golang/go/blob/go1.16.3/" + relpath
										relpath += "#L" + strconv.Itoa(line)
										value_data := []string{pkg, relpath, strconv.Itoa(line), strconv.Itoa(column), name.Name, fmt.Sprintf("\tType: %v", valuespec.Type)}
										gvar_dataset = append(gvar_dataset, value_data)
									}
								}
							}
						}
					}
				}
			}

			// 输出到CSV
			WriteCSVDataset(("../data/func.csv"), func_dataset)
			WriteCSVDataset(("../data/type.csv"), type_dataset)
			WriteCSVDataset(("../data/gvar.csv"), gvar_dataset)
		}
	}
}
