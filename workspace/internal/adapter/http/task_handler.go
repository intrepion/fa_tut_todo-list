package httpadapter

import (
	"net/http"
	"strings"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	service contracts.TaskListService
}

func NewTaskHandler(service contracts.TaskListService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) GetTasks(c echo.Context) error {
	result, err := h.service.ListTasks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TaskHandler) AddTask(c echo.Context) error {
	var request contracts.AddTaskRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
			Message: "invalid request body",
		})
	}

	if strings.TrimSpace(request.Task) == "" {
		return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
			Message: "task must not be blank",
		})
	}

	result, err := h.service.AddTask(request.Task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TaskHandler) RemoveTask(c echo.Context) error {
	taskText := c.QueryParam("task")
	if strings.TrimSpace(taskText) == "" {
		return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
			Message: "task query parameter must not be blank",
		})
	}

	result, err := h.service.RemoveTask(taskText)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, result)
}
