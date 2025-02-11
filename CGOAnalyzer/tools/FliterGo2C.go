package tools

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//IsGo2CDir recursively check if there exists a file(ends with .go) under the given dir contains import "C"
//return true or false
func IsGo2CDir(pth, pkgname string) bool {
	var flag bool
	flag = false
	fset := token.NewFileSet()
	msgReport(os.Stdout, showAll, "now checking %s\n", pth)
	f, _ := parser.ParseDir(fset, pth, nil, parser.ImportsOnly)
	/*if err != nil {
		msgReport(os.Stderr,true,"error occured while parsing %s\n", pth)
		msgReport(os.Stderr,true,err.Error())
	}*/
	for pkg, pkgast := range f {
		if flag {
			break
		}
		msgReport(os.Stdout, showAll, "now parsing:path:%s,pkg:%s\n", pth, pkg)
		for fn, srcfile := range pkgast.Files {
			if flag {
				break
			}
			for _, s := range srcfile.Imports {
				if flag {
					break
				}
				msgReport(os.Stdout, showAll, "pth:%s,pkg:%s,import:%s\n", pth, pkg, s.Path.Value)
				import_name := strings.Trim(s.Path.Value, "\"")
				if import_name == pkgname || strings.HasPrefix(import_name, pkgname+"/") {
					flag = true
					msgReport(os.Stdout, showAll, "hit,go-c found in :%s,pkg:%s\n", fn, pkg)
					break
				}
			}
		}
	}
	return flag
}

func IsValidDir(pth, pkgname string) bool {
	var flag bool
	flag = false
	name := path.Base(pth)
	if strings.HasPrefix(name, ".") {
		return flag
	}
	fs, _ := ioutil.ReadDir(pth)
	for _, fi := range fs {
		if fi.IsDir() {
			if strings.HasPrefix(fi.Name(), `.`) || fi.Name() == `vendor` || fi.Name() == `test` {
				continue
			}
			flag = IsValidDir(filepath.Join(pth, fi.Name()), pkgname)
			if flag {
				break
			}
		}
	}
	if flag {
		return flag
	}
	fset := token.NewFileSet()
	f, err := parser.ParseDir(fset, pth, nil, parser.ImportsOnly)
	if err != nil {
		return false
	}
	for pkg, pkgast := range f {
		if strings.HasSuffix(pkg, "_test") {
			continue
		}
		if flag {
			break
		}
		for _, srcfile := range pkgast.Files {
			if flag {
				break
			}
			for _, s := range srcfile.Imports {
				if flag {
					break
				}
				import_name := strings.Trim(s.Path.Value, "\"")
				if import_name == pkgname || strings.HasPrefix(import_name, pkgname+"/") {
					flag = true
					break
				}
			}
		}
	}
	return flag
}

func IsCGOFile(file *ast.File) bool {
	for _, imp := range file.Imports {
		if imp.Path.Value == `"C"` {
			return true
		}
	}
	return false
}
