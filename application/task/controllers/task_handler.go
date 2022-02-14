package controllers

import (
	"TodoAPI/application/task/entities"
	"TodoAPI/application/task/usecases"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	taskUsecase usecases.TaskUsecaseProtocol
}

func InitTaskHandler(e *echo.Echo, group *echo.Group, taskUsecase usecases.TaskUsecaseProtocol) {
	handler := &TaskHandler{
		taskUsecase: taskUsecase,
	}

	group.GET("/tasks", handler.Fetch)
	group.POST("/tasks", handler.Store)
	group.PUT("/tasks/:id", handler.Update)
	group.DELETE("/tasks/:id", handler.Delete)
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

	ctx := c.Request().Context()

	err := r.taskUsecase.Store(ctx, &task)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, task)

}

func (r *TaskHandler) Update(c echo.Context) error {
	paramId, errInt := strconv.Atoi(c.Param("id"))

	if errInt != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": errInt.Error(),
		})
	}

	isDone, errBool := strconv.ParseBool(c.FormValue("is_done"))

	if errBool != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": errBool.Error(),
		})
	}

	id := int64(paramId)

	task := entities.Task{}
	task.Id = id
	task.Title = c.FormValue("title")
	task.Description = c.FormValue("description")
	task.IsDone = isDone

	ctx := c.Request().Context()

	err := r.taskUsecase.Update(ctx, &task, id)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Updated successfully",
	})
}

func (r *TaskHandler) Delete(c echo.Context) error {

	paramId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": err.Error(),
		})
	}

	id := int64(paramId)

	ctx := c.Request().Context()

	err = r.taskUsecase.Delete(ctx, id)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Delete successfully",
	})
}
