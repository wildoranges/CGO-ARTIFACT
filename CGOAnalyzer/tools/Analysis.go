package tools

import (
	"database/sql"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func BuildRepos(repos []string) {
	db, _ := ConnectSQL()
	var pkgid, repoid, fileID int
	fileID = 175559
	pkgid = 51726
	repoid = 935

	for _, repo := range repos {
		msgReport(os.Stdout, true, "now parsing repo:%s\n", repo)

		// Save Repo Message
		repo_name := filepath.Base(repo)
		info := Repoinfos[RepoIndex[repo_name]]
		data := RepoInfo2String(info)
		data = append(data, strconv.Itoa(repoid))
		err := Insert2Tabel(db, "repository", data)
		if err != nil {
			panic(err)
		}

		alldirs, err := FindAllDirs(repo)
		if err != nil {
			msgReport(os.Stderr, true, err.Error()+"\n")
			return
		}

		for _, repoDir := range alldirs {
			set := token.NewFileSet()
			f, err := parser.ParseDir(set, repoDir, filter, parser.ImportsOnly)
			if err != nil {
				msgReport(os.Stderr, true, err.Error()+"\n")
			}
			for pkg, pkgast := range f {
				if strings.HasSuffix(pkg, `_test`) {
					continue
				}
				// Save Package Message
				pkgData := []string{strconv.Itoa(pkgid), pkg, strconv.Itoa(repoid), repoDir}
				err := Insert2Tabel(db, "package", pkgData)
				if err != nil {
					panic(err)
				}

				msgReport(os.Stdout, showAll, "now parsing pkg:%s\n", pkg)

				// parse file
				for filename, fileast := range pkgast.Files {
					msgReport(os.Stdout, showAll, "now parsing file:%s\n", filename)
					// Save File Message
					fileData := []string{strconv.Itoa(fileID), strconv.Itoa(pkgid), strconv.Itoa(repoid), filename}
					err := Insert2Tabel(db, "file", fileData)
					if err != nil {
						panic(err)
					}

					// Save Import Message
					count := 0
					impDatas := make([]string, 0)
					for _, imp := range fileast.Imports {
						line := set.Position(imp.Pos()).Line
						collumn := set.Position(imp.Pos()).Column
						impData := []string{imp.Path.Value, "0", strconv.Itoa(fileID), strconv.Itoa(line), strconv.Itoa(collumn)}
						impDatas = append(impDatas, impData...)
						count += 1
					}
					if count > 0 {
						err = InsertMuti2Table(db, "import", impDatas, count)
					}
					if err != nil {
						panic(err)
					}
					fileID += 1
				}
				pkgid += 1
			}
		}
		repoid += 1
	}

	db.Close()
}

func UpdateRepo(pkgname string, repoType string, repoID int, db *sql.DB) {
	// 更新Repo类型
	var name string
	if pkgname == "C" {
		name = "cgo"
	} else {
		name = pkgname
	}
	hasPkgType := strings.Contains(repoType, name)
	if !hasPkgType {
		err := UpdateRepoType(db, name, strconv.Itoa(repoID))
		if err != nil {
			panic(err)
		}
	}
}

func filter(fi fs.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), `_test.go`)
}

func AnalyzeRepos(repos []string, pkgNames []string) {
	db, _ := ConnectSQL()
	// var pkgid int

	bar := progressbar.Default(int64(len(repos)))
	for _, repo := range repos {
		msgReport(os.Stdout, true, "now parsing repo:%s\n", repo)

		// Get Repo ID
		repoName := path.Base(repo)
		repoid, repo_type := SelectRepoID(db, repoName)
		for _, pkgName := range pkgNames {
			UpdateRepo(pkgName, repo_type, repoid, db)
		}

		// 获取所有子目录
		alldirs, err := FindAllDirs(repo)
		if err != nil {
			msgReport(os.Stderr, true, err.Error()+"\n")
			return
		}

		// 解析子目录
		for _, repoDir := range alldirs {
			set := token.NewFileSet()
			var f map[string]*ast.Package
			if pkgNames[0] == "C" {
				f, err = parser.ParseDir(set, repoDir, filter, parser.ParseComments)
			} else {
				f, err = parser.ParseDir(set, repoDir, filter, 0)
			}
			if err != nil {
				msgReport(os.Stderr, true, err.Error()+"\n")
			}
			for pkg, pkgast := range f {
				if strings.HasSuffix(pkg, `_test`) {
					continue
				}
				// 开始分析
				msgReport(os.Stdout, showAll, "now parsing pkg:%s\n", pkg)

				for fileName, fileAst := range pkgast.Files {
					fileName = strings.ReplaceAll(fileName, `'`, `''`)
					fileId, err := SelectFileId(db, fileName)
					if err != nil {
						fmt.Println(fileName)
						panic(err)
					}
					msgReport(os.Stdout, showAll, "-- now parsing file:%s\n", fileName)
					for _, pkgname := range pkgNames {
						if pkgname == "C" {
							CntGo2CFunc_File(fileAst, fileId, set, db, fileName)
						} else {
							CntGopkgFunc_File(pkgname, fileAst, fileId, set, db)
						}
					}
				}

			}
		}
		bar.Add(1)
	}

	db.Close()
}
