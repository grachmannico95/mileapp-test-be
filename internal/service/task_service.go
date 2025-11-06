package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/grachmannico95/mileapp-test-be/internal/dto"
	"github.com/grachmannico95/mileapp-test-be/internal/model"
	"github.com/grachmannico95/mileapp-test-be/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService interface {
	Create(ctx context.Context, req dto.CreateTaskRequest) (*model.Task, error)
	GetByID(ctx context.Context, id string) (*model.Task, error)
	List(ctx context.Context, params dto.TaskQueryParams) ([]model.Task, dto.PaginationMeta, error)
	Update(ctx context.Context, id string, req dto.UpdateTaskRequest) (*model.Task, error)
	Delete(ctx context.Context, id string) error
}

type taskServiceImpl struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskServiceImpl{
		taskRepo: taskRepo,
	}
}

func (s *taskServiceImpl) Create(ctx context.Context, req dto.CreateTaskRequest) (*model.Task, error) {
	status := model.TaskStatusPending
	if req.Status != "" {
		status = model.TaskStatus(req.Status)
	}

	priority := model.TaskPriorityMedium
	if req.Priority != "" {
		priority = model.TaskPriority(req.Priority)
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		if req.DueDate.Before(time.Now()) {
			return nil, errors.New("due date must be in the future")
		}
		dueDate = &req.DueDate.Time
	}

	task := model.NewTask(req.Title, req.Description, status, priority, dueDate)

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskServiceImpl) GetByID(ctx context.Context, id string) (*model.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid task ID")
	}

	task, err := s.taskRepo.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, errors.New("task not found")
	}

	return task, nil
}

func (s *taskServiceImpl) List(ctx context.Context, params dto.TaskQueryParams) ([]model.Task, dto.PaginationMeta, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	filters := repository.TaskFilters{
		Status:      params.Status,
		Priority:    params.Priority,
		Search:      params.Search,
		DueDateFrom: params.DueDateFrom,
		DueDateTo:   params.DueDateTo,
		SortBy:      params.SortBy,
		SortOrder:   params.SortOrder,
		Page:        params.Page,
		Limit:       params.Limit,
	}

	tasks, total, err := s.taskRepo.Find(ctx, filters)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	meta := dto.PaginationMeta{
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}

	return tasks, meta, nil
}

func (s *taskServiceImpl) Update(ctx context.Context, id string, req dto.UpdateTaskRequest) (*model.Task, error) {
	task, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		task.Title = req.Title
	}

	if req.Description != "" {
		task.Description = req.Description
	}

	if req.Status != "" {
		task.Status = model.TaskStatus(req.Status)
	}

	if req.Priority != "" {
		task.Priority = model.TaskPriority(req.Priority)
	}

	if req.DueDate != nil && req.DueDate.Time != *task.DueDate {
		if req.DueDate.Before(time.Now()) {
			return nil, errors.New("due date must be in the future")
		}
		task.DueDate = &req.DueDate.Time
	}

	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskServiceImpl) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID")
	}

	return s.taskRepo.Delete(ctx, objectID)
}
