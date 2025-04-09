package models

import "time"

type Task struct {
	ID          int32     `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	StatusID    int32     `db:"status_id"`
	CreatedAt   time.Time `db:"created_at"`
}
