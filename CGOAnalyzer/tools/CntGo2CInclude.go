package tools

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type IncludeItem struct {
	label string
}

const StdHeader int = 0
const UserHeader int = 1
const HeaderUnknown int = 2

var IncludeKind = map[int]string{
	StdHeader:     "StdHeader",
	UserHeader:    "UserHeader",
	HeaderUnknown: "Unknown",
}

func (i IncludeItem) KindText(name string) string {
	if strings.HasPrefix(name, `<`) && strings.HasSuffix(name, `>`) {
		return IncludeKind[StdHeader]
	} else if strings.HasPrefix(name, `"`) && strings.HasSuffix(name, `"`) {
		return IncludeKind[UserHeader]
	}
	return IncludeKind[HeaderUnknown]
}

func (i IncludeItem) Label() string {
	return i.label
}

func CntGo2CLib(repos []string) ItemInfo {

	var Info ItemInfo
	// Info.total = make(map[string]int)
	// Info.local = make(map[string]map[string]int)
	// Info.i = IncludeItem{"header"}
	// Info.pathInfo = make(map[string]StringSet)
	// regStd := regexp.MustCompile(`^.*#include *(<.+\.h>).*$`)
	// regUser := regexp.MustCompile(`^.*#include *(".+\.h").*$`)
	regLib := regexp.MustCompile(`^.*#cgo LDFLAGS: *(-l.+).*$`)
	for _, repo := range repos {

		// curM := make(map[string]int)
		allDirs, err := GetAllDirs(repo)
		if err != nil {
			continue
		}
		for _, subDir := range allDirs {
			set := token.NewFileSet()
			f, err := parser.ParseDir(set, subDir, nil, parser.ParseComments)
			if err != nil {
				msgReport(os.Stderr, true, err.Error()+"\n")
			}
			for pkg, pkgast := range f {
				msgReport(os.Stdout, showAll, "now parsing pkg:%s\n", pkg)
				for filename, srcfile := range pkgast.Files {
					msgReport(os.Stdout, showAll, "now parsing file:%s\n", filename)
					for _, c := range srcfile.Comments {
						cmt := c.Text()
						lines := strings.Split(cmt, "\n")
						for _, line := range lines {
							// headers = regStd.FindStringSubmatch(line)
							// if len(headers) <= 1 {
							// 	headers = regUser.FindStringSubmatch(line)
							// }
							headers := regLib.FindStringSubmatch(line)
							if len(headers) > 1 {
								header := headers[1]
								libs := strings.Split(header, "-l")
								fmt.Printf("found lib ")
								for _, lib := range libs {
									fmt.Printf("%s ", lib)
								}
								fmt.Printf(" in file%s\n", filepath.Base(filename))
								// msgReport(os.Stdout, true, "found lib %s in file %s\n", header, filepath.Base(filename))
								// if Info.pathInfo[header] == nil {
								// 	Info.pathInfo[header] = make(StringSet)
								// }
								// pathRecord(dumpAll, Info, header, filename)
								// if value, ok := curM[header]; ok {
								// 	curM[header] = value + 1
								// } else {
								// 	curM[header] = 1
								// }
							} else {
								continue
							}
						}
					}
				}
			}
		}
		// Info.local[repo] = curM
		// for k, v := range curM {
		// 	if value, ok := Info.total[k]; ok {
		// 		Info.total[k] = value + v
		// 	} else {
		// 		Info.total[k] = v
		// 	}
		// }
	}
	return Info
}
