package models

const (
	StatusCodeTodo = 1
	StatusCodeDone = 2
)

type Status struct {
	Id   int32  `db:"id"`
	Name string `db:"name"`
}
