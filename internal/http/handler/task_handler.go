package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/grachmannico95/mileapp-test-be/internal/dto"
	"github.com/grachmannico95/mileapp-test-be/internal/service"
	"github.com/grachmannico95/mileapp-test-be/internal/util"
)

type TaskHandler struct {
	taskService service.TaskService
	validator   *validator.Validate
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		validator:   validator.New(),
	}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req dto.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("validation failed", util.ParseValidationError(err)...))
		return
	}

	task, err := h.taskService.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	response := dto.ToTaskResponse(task)
	c.JSON(http.StatusCreated, dto.SuccessResponse("task created successfully", response))
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	task, err := h.taskService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
		return
	}

	response := dto.ToTaskResponse(task)
	c.JSON(http.StatusOK, dto.SuccessResponse("task retrieved successfully", response))
}

func (h *TaskHandler) List(c *gin.Context) {
	var params dto.TaskQueryParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("validation failed", util.ParseValidationError(err)...))
		return
	}

	tasks, meta, err := h.taskService.List(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	response := dto.ToTaskListResponse(tasks, meta)
	c.JSON(http.StatusOK, dto.SuccessResponse("tasks retrieved successfully", response))
}

func (h *TaskHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("validation failed", util.ParseValidationError(err)...))
		return
	}

	task, err := h.taskService.Update(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	response := dto.ToTaskResponse(task)
	c.JSON(http.StatusOK, dto.SuccessResponse("task updated successfully", response))
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.taskService.Delete(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("task deleted successfully", nil))
}
