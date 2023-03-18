package main

import (
	"bufio"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

var _init = flag.Bool("init", false, "initialize flag")

func initializeDatabase(db *sqlx.DB) {
	log.Println("initializing database")
	db.MustExec(drop)
	db.MustExec(schema)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		ms := ulid.Now()
		id, _ := ulid.New(ms, r)
		db.MustExec("INSERT INTO tasks (id, text) values (?, ?)", id.String(), "{}")
	}
}

func mainInternal() error {
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
		return err
	}
	if *_init {
		initializeDatabase(db)
		return nil
	}

	q := &TaskQueue{}
	_, err = transaction(db, func(tx *sqlx.Tx) (tasks []*Task, err error) {
		tasks, err = q.getTasks(tx)
		log.Printf("tasks %+v", tasks)
		// トランザクション止めるために入力待ち
		bufio.NewScanner(os.Stdin).Scan()
		return
	})
	return err
}

func main() {
	flag.Parse()
	if err := mainInternal(); err != nil {
		log.Fatal(err)
	}
}
