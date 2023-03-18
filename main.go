package main

import (
	"bufio"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
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
	q := &TaskQueue{}
	_, err = transaction(db, func(tx *sqlx.Tx) (tasks []*Task, err error) {
		tasks, err = q.getTasks(tx)
		log.Printf("tasks %+v", tasks)
		// トランザクション止めるために入力待ち
		bufio.NewScanner(os.Stdin).Scan()
		return
	})
	if err != nil {
		log.Fatal(err)
	}
}
