package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var Connection *sql.DB

func Init() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/insiderDb")
	if err != nil {
		return nil, err
	}
	Connection = db
	return db, nil
}
