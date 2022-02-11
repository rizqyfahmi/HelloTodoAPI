package usecases

import (
	"TodoAPI/app/task/entities"
	"context"
)

type TaskUsecaseProtocol interface {
	Store(c context.Context, task *entities.Task) error
	Update(c context.Context, task *entities.Task, id int64) error
	Fetch(c context.Context) ([]entities.Task, error)
	Delete(c context.Context, id int64) error
}
