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

func (m *DBModel) GetAll() ([]*Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `Select * FROM tasks`
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, err
}

func (m *DBModel) Add(title, description string) error {
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `INSERT INTO tasks (title,description,created_at,updated_at) VALUES ($1,$2,$3,$4)`
	row, err := m.DB.Query(query, title, description, time.Now(), time.Now())
	row.Close()
	if err != nil {
		return err
	}
	return nil
}

func (m *DBModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `DELETE FROM tasks WHERE id=$1`
	row, err := m.DB.QueryContext(ctx, query, id)
	row.Close()
	if err != nil {
		return err
	}
	return nil
}

func (m *DBModel) Update(id int, taskdto TaskDTO) (*Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if t, err := m.Get(id); err != nil || t == nil {
		return nil, err
	}

	query := `UPDATE tasks SET title=$1 , 
	description=$2 , 
	updated_at=$3 WHERE id=$4;`
	row := m.DB.QueryRowContext(ctx, query, taskdto.Title, taskdto.Description, time.Now(), id)

	var task Task

	_ = row.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.UpdatedAt)

	return m.Get(id)
}
