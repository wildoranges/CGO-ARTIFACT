package tools

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
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

var repos_info_path string = `./nosql-repos-info.json`
var RepoIndex map[string]int
var Repoinfos []RepoInfo

func GetRepoInfo() {
	RepoIndex = make(map[string]int)
	b, _ := ioutil.ReadFile(repos_info_path)
	err := json.Unmarshal(b, &Repoinfos)
	if err != nil {
		panic(err)
	}
	for i, repoinfo := range Repoinfos {
		RepoIndex[repoinfo.Repo_name] = i
	}
}

func Bool2String(a bool) string {
	if a {
		return "true"
	} else {
		return "false"
	}
}

func RepoInfo2String(info RepoInfo) []string {
	data := []string{info.Repo_name, info.Url, strconv.Itoa(info.Stars), strconv.Itoa(info.Loc),
		strconv.Itoa(info.Size), strconv.Itoa(info.Forks_count), strconv.Itoa(info.Issues),
		info.Created_at, info.Updated_at, info.Repo_type, info.Pushed_at, info.Description,
		Bool2String(info.Archived), Bool2String(info.Bindings)}
	return data
}
