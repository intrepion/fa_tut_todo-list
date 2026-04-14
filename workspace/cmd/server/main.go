package main

import (
	"log"

	httpadapter "github.com/intrepion/fa_tut_todo-list/workspace/internal/adapter/http"
	storageadapter "github.com/intrepion/fa_tut_todo-list/workspace/internal/adapter/storage"
	"github.com/intrepion/fa_tut_todo-list/workspace/internal/code"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:25616"},
	}))

	store := storageadapter.NewJSONTaskStore("data/tasks.json")
	service := code.NewTaskListService(store)
	handler := httpadapter.NewTaskHandler(service)

	e.GET("/api/tasks", handler.GetTasks)
	e.POST("/api/tasks", handler.AddTask)
	e.DELETE("/api/tasks", handler.RemoveTask)

	log.Fatal(e.Start(":25664"))
}
