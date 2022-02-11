package controllers

import (
	"TodoAPI/app/task/entities"
	"TodoAPI/app/task/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	taskUsecase usecases.TaskUsecaseProtocol
}

func InitTaskHandler(e *echo.Echo, taskUsecase usecases.TaskUsecaseProtocol) {
	handler := &TaskHandler{
		taskUsecase: taskUsecase,
	}

	e.GET("/tasks", handler.Fetch)
	e.POST("/tasks", handler.Store)
}

func (r *TaskHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	result, err := r.taskUsecase.Fetch(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (r *TaskHandler) Store(c echo.Context) error {
	task := entities.Task{}
	task.Title = c.FormValue("title")
	task.Description = c.FormValue("description")

	// err := c.Bind(&task)

	// if err != nil {
	// 	return c.JSON(http.StatusUnprocessableEntity, err.Error())
	// }

	ctx := c.Request().Context()

	err := r.taskUsecase.Store(ctx, &task)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, task)

}
