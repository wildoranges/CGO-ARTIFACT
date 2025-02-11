package tools

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

var insertPattern map[string]string = map[string]string{
	"repository":        "INSERT INTO repository(repo_name, url, stars, loc, size, forks_count, issues, created_at, updated_at, repo_type, pushed_at, description, archived, bindings, id) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)",
	"package":           "INSERT INTO package(id, name, repo_id, path) VALUES($1, $2, $3, $4)",
	"file":              "INSERT INTO file(id, pkg_id, repo_id, path) VALUES($1, $2, $3, $4)",
	"cgo_function":      "INSERT INTO cgo_function(line, file_id, func_name, func_description, collumn) VALUES($1, $2, $3, $4, $5)",
	"cgo_header":        "INSERT INTO cgo_header(line, file_id, header_name, type, collumn) VALUES($1, $2, $3, $4, $5)",
	"cgo_lib":           "INSERT INTO cgo_lib(line, file_id, lib_name, lib_description, collumn) VALUES($1, $2, $3, $4, $5)",
	"crypto_invocation": "INSERT INTO crypto_invocation(line, file_id, func_name, pkgname, collumn) VALUES($1, $2, $3, $4, $5)",
	"math_invocation":   "INSERT INTO math_invocation(line, file_id, func_name, pkgname, collumn) VALUES($1, $2, $3, $4, $5)",
	"import":            "INSERT INTO import(name, type, file_id, line, collumn) VALUES($1, $2, $3, $4, $5)",
	"export_function":   "INSERT INTO export_function(line, file_id, func_name, collumn) VALUES($1, $2, $3, $4)",
	"net_invocation":    "INSERT INTO net_invocation(line, file_id, func_name, pkgname, collumn) VALUES($1, $2, $3, $4, $5)",
}

func ConnectSQL() (*sql.DB, error) {
	connStr := "postgres://postgres:s4plususer@localhost/cgo?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Insert2Tabel(db *sql.DB, table string, data []string) error {
	stmt, err := db.Prepare(insertPattern[table])
	if err != nil {
		return err
	}
	y := make([]interface{}, len(data))
	for i, d := range data {
		y[i] = d
	}
	_, err = stmt.Exec(y...)
	// id, err1 := result.LastInsertId()
	// fmt.Println(id, err1)
	return err
}

func InsertMuti2Table(db *sql.DB, table string, data []string, num int) error {
	ptn := insertPattern[table]
	n0 := strings.Count(ptn, "$")
	n := n0 + 1

	for i := 1; i < num; i++ {
		ptn += ", ("
		for i := 0; i < n0-1; i++ {
			s := fmt.Sprintf("$%v,", n)
			n += 1
			ptn += s
		}
		s := fmt.Sprintf("$%v)", n)
		n += 1
		ptn += s
	}

	stmt, err := db.Prepare(ptn)
	if err != nil {
		return err
	}
	y := make([]interface{}, len(data))
	for i, d := range data {
		y[i] = d
	}
	_, err = stmt.Exec(y...)
	return err
}

func UpdateRepoType(db *sql.DB, repo_type string, id string) error {
	stmt, err := db.Prepare("update repository set repo_type=CONCAT_WS(';', repo_type, $1::text) where id=$2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(repo_type, id)
	return err
}

func SelectRepoID(db *sql.DB, repoName string) (id int, repo_type string) {
	row := db.QueryRow(`SELECT id, repo_type FROM repository where repo_name='` + repoName + `'`)
	row.Scan(&id, &repo_type)
	return
}

func SelectPkgId(db *sql.DB, pkgname, path string) (id int, err error) {
	row := db.QueryRow(`SELECT id FROM package where pkg_path='` + path + `' and pkg_name='` + pkgname + `'`)
	err = row.Scan(&id)
	return
}

func SelectFileId(db *sql.DB, path string) (id int, err error) {
	row := db.QueryRow(`SELECT id FROM file where path='` + path + `'`)
	err = row.Scan(&id)
	return
}
