package tools

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type data struct {
	pkg           string
	file          string
	type_name     string
	function_name string
}

type MethodVisitor struct {
	ExportedTypes map[string]int
	dataset       *[]data
	set           *token.FileSet
	pkg           string
	tag           int
}

func (m MethodVisitor) Visit(n ast.Node) ast.Visitor {
	if m.tag == 1 {
		if typespec, ok := n.(*ast.TypeSpec); ok {
			if isExported(typespec.Name.Name) {
				m.ExportedTypes[typespec.Name.Name] = 0
				if it, ok := typespec.Type.(*ast.InterfaceType); ok {
					for _, method := range it.Methods.List {
						for _, name := range method.Names {
							line := m.set.Position(method.Pos()).Line
							filename := m.set.File(method.Pos())
							relpath, _ := filepath.Rel("/home/dby/go-c/new-go/go", filename.Name())
							relpath = "https://github.com/golang/go/blob/go1.16.3/" + relpath
							relpath += "#L" + strconv.Itoa(line)
							new_data := data{pkg: m.pkg, file: relpath, type_name: typespec.Name.Name, function_name: name.Name}
							*m.dataset = append(*m.dataset, new_data)
						}
					}
				}
			}
		}
	} else {
		if function, ok := n.(*ast.FuncDecl); ok {
			if isExported(function.Name.Name) && function.Recv != nil {
				for _, item := range function.Recv.List {
					// for _, name := range item.Names {
					// 	fmt.Printf("name.Name: %v\n", name.Name)
					// }
					if star, ok := item.Type.(*ast.StarExpr); ok {
						if ident, ok := star.X.(*ast.Ident); ok {
							if _, ok := m.ExportedTypes[ident.Name]; ok {
								line := m.set.Position(function.Pos()).Line
								filename := m.set.File(function.Pos())
								relpath, _ := filepath.Rel("/home/dby/go-c/new-go/go", filename.Name())
								relpath = "https://github.com/golang/go/blob/go1.16.3/" + relpath
								relpath += "#L" + strconv.Itoa(line)
								new_data := data{pkg: m.pkg, file: relpath, type_name: ident.Name, function_name: function.Name.Name}
								*m.dataset = append(*m.dataset, new_data)
							}
						}
					}
					if ident, ok := item.Type.(*ast.Ident); ok {
						if _, ok := m.ExportedTypes[ident.Name]; ok {
							line := m.set.Position(function.Pos()).Line
							filename := m.set.File(function.Pos())
							relpath, _ := filepath.Rel("/home/dby/go-c/new-go/go", filename.Name())
							relpath = "https://github.com/golang/go/blob/go1.16.3/" + relpath
							relpath += "#L" + strconv.Itoa(line)
							new_data := data{pkg: m.pkg, file: relpath, type_name: ident.Name, function_name: function.Name.Name}
							*m.dataset = append(*m.dataset, new_data)
						}
					}
				}
			}
		}
	}
	return m
}

func GetMethod(dir string) {
	header := []string{"package", "file", "type name", "method name"}
	if _, err := os.Stat("../data/method.csv"); os.IsNotExist(err) {
		WriteCSVData("../data/method.csv", header)
	}

	var dataset [][]string = make([][]string, 0)

	set := token.NewFileSet()
	f, err := parser.ParseDir(set, dir, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	for pkg, pkgAST := range f {
		if strings.HasSuffix(pkg, "_test") {
			continue
		}
		for filename, fileAST := range pkgAST.Files {
			dataset = make([][]string, 0)
			if strings.HasSuffix(filename, "_test.go") {
				continue
			}
			temp := make([]data, 0)
			m := MethodVisitor{make(map[string]int), &temp, set, pkg, 1}
			ast.Walk(m, fileAST)
			m.tag = 2
			ast.Walk(m, fileAST)
			for _, d := range *m.dataset {
				new_data := make([]string, 0)
				new_data = append(new_data, d.pkg)
				new_data = append(new_data, d.file)
				new_data = append(new_data, d.type_name)
				new_data = append(new_data, d.function_name)
				dataset = append(dataset, new_data)
			}
			WriteCSVDataset("../data/method.csv", dataset)
		}
	}
}
