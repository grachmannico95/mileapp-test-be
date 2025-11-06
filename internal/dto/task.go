package dto

import (
	"time"

	"github.com/grachmannico95/mileapp-test-be/internal/model"
)

type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required,min=3,max=200"`
	Description string    `json:"description" binding:"max=2000"`
	Status      string    `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
	Priority    string    `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *JSONTime `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string    `json:"title" binding:"required,min=3,max=200"`
	Description string    `json:"description" binding:"max=2000"`
	Status      string    `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
	Priority    string    `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *JSONTime `json:"due_date"`
}

type TaskQueryParams struct {
	Page        int    `form:"page" binding:"omitempty,min=1"`
	Limit       int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Status      string `form:"status" binding:"omitempty,oneof=pending in_progress completed"`
	Priority    string `form:"priority" binding:"omitempty,oneof=low medium high"`
	Search      string `form:"search"`
	DueDateFrom string `form:"due_date_from"`
	DueDateTo   string `form:"due_date_to"`
	SortBy      string `form:"sort_by" binding:"omitempty,oneof=created_at updated_at due_date priority title"`
	SortOrder   string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

type TaskResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

type TaskListResponse struct {
	Tasks []TaskResponse `json:"tasks"`
	Meta  PaginationMeta `json:"meta"`
}

type JSONTime struct {
	time.Time
}

func (jt *JSONTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" || s == `""` {
		return nil
	}

	// Remove quotes
	s = s[1 : len(s)-1]

	// Parse ISO8601 format
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	jt.Time = t
	return nil
}

func ToTaskResponse(task *model.Task) TaskResponse {
	return TaskResponse{
		ID:          task.ID.Hex(),
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		Priority:    string(task.Priority),
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToTaskListResponse(tasks []model.Task, meta PaginationMeta) TaskListResponse {
	taskResponses := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = ToTaskResponse(&task)
	}

	return TaskListResponse{
		Tasks: taskResponses,
		Meta:  meta,
	}
}
