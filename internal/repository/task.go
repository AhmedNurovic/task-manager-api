package repository

import (
	"context"
	"database/sql"

	"github.com/ahmednurovic/task-manager-api/internal/model"
	"github.com/jmoiron/sqlx"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetAllForUser(ctx context.Context, userID int64) ([]*model.Task, error)
	GetByID(ctx context.Context, taskID int64) (*model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, taskID int64, userID int64) error
}

type TaskRepositoryImpl struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepositoryImpl {
	return &TaskRepositoryImpl{db: db}
}

func (r *TaskRepositoryImpl) Create(ctx context.Context, task *model.Task) error {
	query := `INSERT INTO tasks (user_id, title, status) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, task.UserID, task.Title, task.Status).Scan(&task.ID)
}

func (r *TaskRepositoryImpl) GetAllForUser(ctx context.Context, userID int64) ([]*model.Task, error) {
	var tasks []*model.Task
	query := `SELECT id, user_id, title, status FROM tasks WHERE user_id = $1`
	err := r.db.SelectContext(ctx, &tasks, query, userID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepositoryImpl) GetByID(ctx context.Context, taskID int64) (*model.Task, error) {
	var task model.Task
	query := `SELECT id, user_id, title, status FROM tasks WHERE id = $1`
	err := r.db.GetContext(ctx, &task, query, taskID)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepositoryImpl) Update(ctx context.Context, task *model.Task) error {
	query := `UPDATE tasks SET title = $1, status = $2 WHERE id = $3 AND user_id = $4`
	result, err := r.db.ExecContext(ctx, query, task.Title, task.Status, task.ID, task.UserID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TaskRepositoryImpl) Delete(ctx context.Context, taskID int64, userID int64) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, taskID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
