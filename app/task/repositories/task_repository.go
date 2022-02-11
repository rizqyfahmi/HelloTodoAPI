package repositories

import (
	"TodoAPI/app/task/entities"
	"context"
	"database/sql"
	"fmt"
)

type TaskRepository struct {
	db *sql.DB
}

func InitTaskRepository(db *sql.DB) TaskRepositoryProtocol {
	return &TaskRepository{
		db: db,
	}
}

func (repo *TaskRepository) Store(ctx context.Context, task *entities.Task) error {

	query := `INSERT INTO task(id, title, description, is_done, updated_at, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	stmt, errPrepare := repo.db.PrepareContext(ctx, query)

	if errPrepare != nil {
		return errPrepare
	}

	res, errExec := stmt.ExecContext(ctx, task.Id, task.Title, task.Description, task.IsDone, task.UpdatedAt, task.CreatedAt)
	if errExec != nil {
		return errExec
	}

	_, errAffected := res.LastInsertId()
	if errAffected != nil {
		return errAffected
	}

	return nil
}

func (repo *TaskRepository) Update(ctx context.Context, task *entities.Task) error {

	query := `UPDATE task SET title=?, description=?, is_done=?, updated_at=? WHERE id=?`
	stmt, errPrepare := repo.db.PrepareContext(ctx, query)
	if errPrepare != nil {
		return errPrepare
	}

	res, errExec := stmt.ExecContext(ctx, task.Title, task.Description, task.IsDone, task.UpdatedAt, task.Id)
	if errExec != nil {
		return errExec
	}

	row, errAffected := res.RowsAffected()
	if errAffected != nil {
		return errAffected
	}

	if row != 1 {
		return fmt.Errorf("weird behavior. total affected: %d", row)
	}

	return nil

}

func (repo *TaskRepository) FindByID(ctx context.Context, id int64) (entities.Task, error) {

	task := entities.Task{}

	query := `SELECT id, title, description, updated_at, created_at FROM task WHERE task.id = ?`
	err := repo.db.QueryRowContext(ctx, query, id).Scan(&task.Id, &task.Title, &task.Description, &task.UpdatedAt, &task.CreatedAt)

	if err != nil {
		return entities.Task{}, err
	}

	return task, nil

}

func (repo *TaskRepository) Fetch(ctx context.Context) ([]entities.Task, error) {
	query := `SELECT id, title, description, is_done, updated_at, created_at FROM task`
	rows, err := repo.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	tasks := make([]entities.Task, 0)
	for rows.Next() {
		task := entities.Task{}
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.IsDone, &task.UpdatedAt, &task.CreatedAt)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (repo *TaskRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM task WHERE task.id = ?`
	stmt, errPrepare := repo.db.PrepareContext(ctx, query)
	if errPrepare != nil {
		return errPrepare
	}

	res, errExec := stmt.ExecContext(ctx, id)
	if errExec != nil {
		return errExec
	}

	row, errAffected := res.RowsAffected()
	if errAffected != nil {
		return errAffected
	}

	if row != 1 {
		return fmt.Errorf("weird behavior. total affected: %d", row)
	}

	return nil

}
