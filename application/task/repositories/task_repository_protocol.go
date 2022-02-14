package repositories

import (
	"TodoAPI/application/task/entities"
	"context"
)

type TaskRepositoryProtocol interface {
	Store(ctx context.Context, task *entities.Task) error
	Update(ctx context.Context, task *entities.Task) error
	FindByID(ctx context.Context, id int64) (entities.Task, error)
	Fetch(ctx context.Context) ([]entities.Task, error)
	Delete(ctx context.Context, id int64) error
}
