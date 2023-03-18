package main

import (
	"time"
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
