package database

import (
	"database/sql"
	"fmt"
	"log"
)
import _ "github.com/go-sql-driver/mysql"

var Db *sql.DB

func init() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=30s", "root", "foxconn", "193.112.35.124:3306", "fdl"))
	if err != nil {
		log.Panic(err)
	}
	Db = db
}
