package main

import (
	delivery "TodoAPI/app/task/controllers"
	"TodoAPI/app/task/repositories"
	"TodoAPI/app/task/usecases"
	"TodoAPI/middlewares"
	Mysql "TodoAPI/modules/mysql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {

	db := Mysql.Connect()
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	e := echo.New()

	middleware := middlewares.InitMiddleware()
	repo := repositories.InitTaskRepository(db)
	usecase := usecases.InitTaskUsecase(repo, 2*time.Minute)
	delivery.InitTaskHandler(e, usecase)

	e.Use(middleware.HandleCORS)
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello World",
		})
	})

	e.Start(":8081")

}
