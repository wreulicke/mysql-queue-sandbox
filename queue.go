package main

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type TaskQueue struct {
}

func transaction[T any](db *sqlx.DB, f func(tx *sqlx.Tx) (T, error)) (T, error) {
	var t T
	tx, err := db.Beginx()
	if err != nil {
		return t, err
	}
	t, err = f(tx)
	if err != nil {
		tx.Rollback()
		return t, err
	}
	return t, tx.Commit()
}

func (q *TaskQueue) getTasks(tx *sqlx.Tx) ([]*Task, error) {
	log.Println("start query")
	var tasks []*Task
	err := tx.Select(&tasks, `
	SELECT * FROM tasks ORDER BY id LIMIT 1 FOR UPDATE SKIP LOCKED;
`)
	log.Println("end query")
	if err != nil {
		return nil, err
	}
	return tasks, err
}
