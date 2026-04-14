package httpadapter

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/intrepion/fa_tut_todo-list/workspace/internal/contracts"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	service contracts.TaskService
}

func NewTaskHandler(service contracts.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) ListTasks(c echo.Context) error {
	result, err := h.service.ListTasks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TaskHandler) CreateTask(c echo.Context) error {
	var request contracts.CreateTaskRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
			Message: "invalid request body",
		})
	}

	task, err := h.service.CreateTask(request.Text)
	if err != nil {
		if errors.Is(err, contracts.ErrTaskTextBlank) {
			return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(c echo.Context) error {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
			Message: "task id must be an integer",
		})
	}

	task, err := h.service.GetTask(taskID)
	if err != nil {
		if errors.Is(err, contracts.ErrTaskNotFound) {
			return c.JSON(http.StatusNotFound, contracts.ErrorResponse{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c echo.Context) error {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, contracts.ErrorResponse{
			Message: "task id must be an integer",
		})
	}

	err = h.service.DeleteTask(taskID)
	if err != nil {
		if errors.Is(err, contracts.ErrTaskNotFound) {
			return c.JSON(http.StatusNotFound, contracts.ErrorResponse{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{
			Message: "internal server error",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
