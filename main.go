package main

import (
	AuthDelivery "TodoAPI/app/auth/controllers"
	AuthRepository "TodoAPI/app/auth/repositories"
	AuthUsecase "TodoAPI/app/auth/usecases"
	TaskDelivery "TodoAPI/app/task/controllers"
	TaskRepository "TodoAPI/app/task/repositories"
	TaskUsecase "TodoAPI/app/task/usecases"
	TodoMiddleware "TodoAPI/middlewares"
	Mysql "TodoAPI/modules/mysql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type Request struct {
	Id    string `json:"id" form:"id" query:"id"`
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

func main() {

	db := Mysql.Connect()
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	e := echo.New()
	const timeout = 2 * time.Minute

	mw := TodoMiddleware.InitMiddleware()
	e.Use(mw.HandleCORS)

	restricted := e.Group("/restricted")
	restricted.Use(mw.VerifyToken)
	// TASK MODULE
	taskRepo := TaskRepository.InitTaskRepository(db)
	taskUsecase := TaskUsecase.InitTaskUsecase(taskRepo, timeout)
	TaskDelivery.InitTaskHandler(e, restricted, taskUsecase)
	// AUTH MODULE
	authRepo := AuthRepository.InitAuthRepository(db)
	authUsecase := AuthUsecase.InitAuthUsecase(authRepo, timeout)
	AuthDelivery.InitAuthHandler(e, authUsecase)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello World",
		})
	})

	e.POST("/bind", func(c echo.Context) error {
		u := new(Request)
		if err := c.Bind(u); err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello World",
			"data":    u,
		})
	})

	e.Start(":8081")

}
