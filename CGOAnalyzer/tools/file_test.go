package tools_test

import (
	"anatool/tools"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

type RepoInfo struct {
	Repo_name   string
	Url         string
	Stars       int
	Loc         int
	Size        int
	Forks_count int
	Issues      int
	Created_at  string
	Updated_at  string
	Pushed_at   string
	Description string
	Archived    bool
	Repo_type   string
	Bindings    bool
}

func Test_Name(t *testing.T) {
	fileinfo, _ := ioutil.ReadDir("/data/github_go/crypto-repos/")
	for _, si := range fileinfo {
		if si.IsDir() {
			fmt.Println(si.Name())
		}
	}
}

func TestJson(t *testing.T) {
	// repoinfos := make(map[RepoInfo]void)
	var repoinfos []RepoInfo
	// repoinfos := make([]RepoInfo, 900)
	// b := `{"repo_name": "kubefwd", "url": "https://github.com/txn2/kubefwd", "stars": 2624, "loc": 1989, "size": 16264, "forks_count": 137, "issues": 7, "created_at": "2018-08-05 22:05:58", "updated_at": "2021-10-01 15:37:10", "repo_type": ""}`
	b, _ := ioutil.ReadFile("../repos-info.json")
	err := json.Unmarshal(b, &repoinfos)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(len(repoinfos))
	fmt.Println(repoinfos[10].Archived)
	fmt.Println(repoinfos[10].Bindings)
}

func TestAllDirs(t *testing.T) {
	dirs, _ := tools.GetAllDirs(`/home/dby/go-c/mygo/go117`)
	dirs2, _ := tools.FindAllDirs(`/home/dby/go-c/mygo/go117`)
	f1Set := make(tools.StringSet)
	for _, dir := range dirs {
		f1Set.Insert(dir)
	}
	for _, dir := range dirs2 {
		if _, ok := f1Set[dir]; !ok {
			fmt.Println(dir)
		}
	}
}
