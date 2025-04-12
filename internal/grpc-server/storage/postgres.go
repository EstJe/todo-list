package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/EstJe/todo-list/app/domain/models"
	"github.com/EstJe/todo-list/app/internal/lib/op"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db}
}

func (p *Postgres) CreateTask(ctx context.Context, title string, description string) (int32, error) {
	query := `INSERT INTO tasks (title, description) VALUES ($1, $2)RETURNING id`

	var id int32
	err := p.db.QueryRowContext(ctx, query, title, description).Scan(&id)
	if err != nil {
		return 0, op.Wrap(err)
	}

	return id, nil
}

func (p *Postgres) DeleteTask(ctx context.Context, id int32) error {
	query := "DELETE FROM tasks WHERE id = $1"

	res, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return op.Wrap(err)
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return op.Wrap(err)
	}

	if cnt == 0 {
		return op.Wrap(ErrTaskNotFound)
	}

	return nil
}

func (p *Postgres) DoneTask(ctx context.Context, id int32) error {
	query := fmt.Sprintf("UPDATE tasks SET status_id=%d WHERE id = $1", models.StatusCodeDone)

	res, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return op.Wrap(err)
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return op.Wrap(err)
	}

	if cnt == 0 {
		return op.Wrap(ErrTaskNotFound)
	}

	return nil
}

func (p *Postgres) Tasks(ctx context.Context) ([]models.Task, error) {
	query := "SELECT id, title, description, status_id, created_at FROM tasks"

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, op.Wrap(err)
	}
	defer rows.Close()

	tasks := make([]models.Task, 0, 1)
	for rows.Next() {
		t := models.Task{}

		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.StatusID, &t.CreatedAt)
		if err != nil {
			return nil, op.Wrap(err)
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, op.Wrap(err)
	}

	return tasks, nil
}
