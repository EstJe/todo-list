package todogrpc

import tdlapi "github.com/EstJe/todo-list/internal/protos/todo"

type serverAPI struct {
	tdlapi.UnimplementedTodoListServer
}
