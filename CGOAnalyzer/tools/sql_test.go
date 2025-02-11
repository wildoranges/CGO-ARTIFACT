package tools

import (
	"fmt"
	"strings"
	"testing"
)

func TestConnectSQL(t *testing.T) {
	_, err := ConnectSQL()
	if err != nil {
		t.Error(err)
	}
}

func TestInsert2Tabel(t *testing.T) {
	db, err := ConnectSQL()
	if err != nil {
		t.Error(err)
	}
	d := []string{"1", "1", "TESTNET", "net", "2"}
	err = Insert2Tabel(db, "net_invocation", d)
	if err != nil {
		t.Error(err)
	}
}

func TestSelectRepoID(t *testing.T) {
	db, err := ConnectSQL()
	if err != nil {
		t.Error(err)
	}
	repoName := "go"
	id, repo_type := SelectRepoID(db, repoName)
	fmt.Println(id, repo_type)
}

func TestUpdateRepoType(t *testing.T) {
	db, err := ConnectSQL()
	if err != nil {
		t.Error(err)
	}
	err = UpdateRepoType(db, "crypto", "go")
	if err != nil {
		t.Error(err)
	}
	db.Close()
}

func TestSelectPkgId(t *testing.T) {
	db, err := ConnectSQL()
	if err != nil {
		t.Error(err)
	}
	id, err := SelectPkgId(db, "main", "/data/github_go/repos/under-the-hood/demo/ch05boot")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(id)
	db.Close()
}

func TestInsertMuti2Table(t *testing.T) {
	db, _ := ConnectSQL()
	data := []string{"p1", "0", "1", "1", "1", "p2", "0", "1", "1", "1"}
	err := InsertMuti2Table(db, "import", data, 2)
	if err != nil {
		t.Error(err)
	}
	db.Close()
}

func TestSelectFileID(t *testing.T) {
	db, _ := ConnectSQL()
	name := `/data/github_go/all-repos/LeetCode-Go/leetcode/0118.Pascals-Triangle/118. Pascal's Triangle.go`
	name = strings.ReplaceAll(name, `'`, `''`)
	id, err := SelectFileId(db, name)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(id)
	db.Close()
}
