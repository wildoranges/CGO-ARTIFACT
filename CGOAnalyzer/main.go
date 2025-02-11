package main

import (
	"anatool/tools"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// 获取所有包含对pkgname引用的目录
func FilterVallidDir(dirpath, pkgname string) []string {
	fileinfo, err := ioutil.ReadDir(dirpath)
	fmt.Printf("now getting valid repos(%s) from %s\n", pkgname, dirpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error occured while reading dir:%s:\n", dirpath)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	set := make(tools.StringSet)
	for _, si := range fileinfo {
		if _, ok := tools.RepoIndex[si.Name()]; ok {
			// if si.IsDir() {
			curRepo := si.Name()
			fullpath := filepath.Join(dirpath, curRepo)
			isvalid := tools.IsValidDir(fullpath, pkgname)
			if isvalid {
				set[fullpath] = tools.VoidMem
				continue
			}
		}
		// }
	}

	var Validdirs []string
	fmt.Printf("valid repos:\n")
	for p := range set {
		Validdirs = append(Validdirs, p)
		fmt.Println(p)
	}
	fmt.Println("valid repos num: " + strconv.Itoa(len(Validdirs)))

	return Validdirs
}

// func CGO(dirpath string, topx *int, core *int) {
// 	Validdirs := FilterVallidDir(dirpath, "C")

// 	fmt.Printf("now parsing valid repos ...\n")
// 	tools.CntGo2CFunc(Validdirs)
// 	fmt.Printf("finish parsing!\n")
// }

// func GO(dirpath string) {
// 	err := os.RemoveAll("../data")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = os.Mkdir("../data", 0777)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fileinfo, err := ioutil.ReadDir(dirpath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, fi := range fileinfo {
// 		curpath := filepath.Join(dirpath, fi.Name())
// 		if fi.IsDir() {
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			tools.CntGoFunc(curpath)
// 			tools.GetMethod(curpath)
// 		}
// 	}
// }

// func GoPkg(pkgname string, dirpath string) {
// 	Validdirs := FilterVallidDir(dirpath, pkgname)

// 	fmt.Printf("now parsing valid repos ...\n")
// 	tools.CntGopkgFunc(Validdirs, pkgname)
// 	fmt.Printf("finish parsing!\n")
// }

func main() {
	start := time.Now()
	pathPtr := flag.String(`path`, `./`, `dir path of the repos`)
	showAll := flag.Bool(`showall`, false, `show output details`)
	// pkgname := flag.String(`pkgname`, `C`, `Package to be anlyzed`)
	build := flag.Bool(`build`, false, `True: Build repositories and packages; False: Implement analysis`)
	flag.Parse()
	dirpath := *pathPtr
	tools.SetFlags(*showAll)

	pkgNames := []string{"C", "crypto", "net", "math"}
	// pkgNames := "crypto"

	tools.GetRepoInfo()

	fmt.Println(pkgNames, len(pkgNames))
	if *build {
		fmt.Println("Now start to get repo info...")
		var validdir []string
		fileinfo, err := ioutil.ReadDir(dirpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error occured while reading dir:%s:\n", dirpath)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
		for _, si := range fileinfo {
			if _, ok := tools.RepoIndex[si.Name()]; ok {
				p := filepath.Join(dirpath, si.Name())
				validdir = append(validdir, p)
			}
		}
		fmt.Println("Now start to build...")
		tools.BuildRepos(validdir)
	} else {
		validDir := make([]string, 0)
		validDirSet := make(tools.StringSet)
		for _, pkgName := range pkgNames {
			validdir := FilterVallidDir(dirpath, pkgName)
			for _, dir := range validdir {
				if _, ok := validDirSet[dir]; !ok {
					validDirSet.Insert(dir)
				}
			}
		}
		for dir := range validDirSet {
			validDir = append(validDir, dir)
		}
		fmt.Println("repo nums: ", len(validDir))
		tools.AnalyzeRepos(validDir, pkgNames)
	}
	fmt.Println("time duration:", time.Since(start))
}

// func main() {
// 	pkgNames := []string{"crypto", "math"}
// 	validDir := []string{"/data/github_go/all-repos/LeetCode-Go"}
// 	tools.AnalyzeRepos(validDir, pkgNames)
// }
