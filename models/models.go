package models

import (
	"context"
	"database/sql"
	"time"
)

type Models struct {
	DB DBModel
}

type DBModel struct {
	DB *sql.DB
}

func NewModels(db *sql.DB) Models {
	return Models{DB: DBModel{
		DB: db,
	}}
}

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

func (m *DBModel) Get(id int) (*Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `Select * FROM tasks WHERE id=$1`
	row := m.DB.QueryRowContext(ctx, query, id)

	var task Task

	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &task, err
}
