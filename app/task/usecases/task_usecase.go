package usecases

import (
	"TodoAPI/app/task/entities"
	"TodoAPI/app/task/repositories"
	"context"
	"math/rand"
	"time"
)

type TaskUsecase struct {
	taskRepo repositories.TaskRepositoryProtocol
	timeout  time.Duration
}

func InitTaskUsecase(taskRepo repositories.TaskRepositoryProtocol, timeout time.Duration) TaskUsecaseProtocol {
	return &TaskUsecase{
		taskRepo: taskRepo,
		timeout:  timeout,
	}
}

func (usecase *TaskUsecase) Store(c context.Context, task *entities.Task) error {

	ctx, cancel := context.WithTimeout(c, usecase.timeout)
	defer cancel()

	const min int64 = 1000
	const max int64 = 9999

	task.Id = rand.Int63n(max-min) + min
	task.IsDone = false
	task.UpdatedAt = time.Now()
	task.CreatedAt = time.Now()

	err := usecase.taskRepo.Store(ctx, task)

	if err != nil {
		return err
	}

	return nil
}

func (usecase *TaskUsecase) Update(c context.Context, task *entities.Task, id int64) error {

	ctx, cancel := context.WithTimeout(c, usecase.timeout)
	defer cancel()

	if _, errFind := usecase.taskRepo.FindByID(ctx, id); errFind != nil {
		return errFind
	}

	task.UpdatedAt = time.Now()

	err := usecase.taskRepo.Update(ctx, task)

	if err != nil {
		return err
	}

	return nil
}

func (usecase *TaskUsecase) Fetch(c context.Context) ([]entities.Task, error) {
	ctx, cancel := context.WithTimeout(c, usecase.timeout)
	defer cancel()

	result, err := usecase.taskRepo.Fetch(ctx)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (usecase *TaskUsecase) Delete(c context.Context, id int64) error {

	ctx, cancel := context.WithTimeout(c, usecase.timeout)
	defer cancel()

	err := usecase.taskRepo.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil

}
