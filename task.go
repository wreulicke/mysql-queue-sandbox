package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

// テストのため、ここでスキーマ定義している
var drop = `
DROP TABLE IF EXISTS tasks;`

var schema = `
CREATE TABLE IF NOT EXISTS tasks (
		id char(26) PRIMARY KEY,
    text text,
		done boolean DEFAULT false,
		created_at datetime DEFAULT CURRENT_TIMESTAMP
);`

type Task struct {
	ID        string
	Text      string
	Done      bool
	CreatedAt time.Time `db:"created_at"`
}

// テストのため、ここでスキーマ定義している
func init() {
	c := mysql.Config{
		User:      "root",
		Passwd:    "password",
		DBName:    "queue-test",
		Net:       "tcp",
		Addr:      "localhost:3306",
		Collation: "utf8mb4_unicode_ci",
		ParseTime: true,
		Loc:       time.UTC,
	}
	db, err := sqlx.Open("mysql", c.FormatDSN())
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(drop)
	db.MustExec(schema)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		ms := ulid.Now()
		id, _ := ulid.New(ms, r)
		db.MustExec("INSERT INTO tasks (id, text) values (?, ?)", id.String(), "{}")
	}
}
