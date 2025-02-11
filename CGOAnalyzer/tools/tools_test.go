package tools

import (
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
)

func DoAnalysis(repo string) {
	msgReport(os.Stdout, showAll, "now parsing repo:%s\n", repo)

	// 获取所有子目录
	alldirs, err := FindAllDirs(repo)
	if err != nil {
		msgReport(os.Stderr, true, err.Error()+"\n")
		return
	}

	// 解析子目录
	for _, repoDir := range alldirs {
		set := token.NewFileSet()
		f, err := parser.ParseDir(set, repoDir, nil, parser.ParseComments)
		if err != nil {
			msgReport(os.Stderr, true, err.Error()+"\n")
		}
		for pkg, pkgast := range f {
			if strings.HasSuffix(pkg, `_test`) {
				continue
			}
			// 开始分析
			msgReport(os.Stdout, showAll, "now parsing pkg:%s\n", pkg)
			// CntC2GoFunc_Pkg(pkgast, 1, set, nil)
			for _, file := range pkgast.Files {
				CntImport(file)
			}
		}
	}
}

func TestC2GoFunc(t *testing.T) {
	DoAnalysis("/home/dby/go-c/go/")
}

func TestImport(t *testing.T) {
	DoAnalysis("/data/github_go/all-repos/wtf/")
}
