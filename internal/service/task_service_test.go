package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/grachmannico95/mileapp-test-be/internal/dto"
	"github.com/grachmannico95/mileapp-test-be/internal/model"
	"github.com/grachmannico95/mileapp-test-be/internal/repository"
	"github.com/grachmannico95/mileapp-test-be/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTaskService_Create_Success(t *testing.T) {
	// Setup
	mockTaskRepo := mocks.NewMockTaskRepository(t)
	taskService := NewTaskService(mockTaskRepo)

	// Test data
	futureDate := time.Now().Add(24 * time.Hour)
	req := dto.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
		Priority:    "high",
		DueDate:     &dto.JSONTime{Time: futureDate},
	}

	// Mock expectations
	mockTaskRepo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(task *model.Task) bool {
			return task.Title == req.Title &&
				task.Description == req.Description &&
				task.Status == model.TaskStatusPending &&
				task.Priority == model.TaskPriorityHigh
		})).
		Return(nil).
		Once()

	// Execute
	task, err := taskService.Create(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, req.Title, task.Title)
	assert.Equal(t, req.Description, task.Description)
}

func TestTaskService_Create_PastDueDate(t *testing.T) {
	// Setup
	mockTaskRepo := mocks.NewMockTaskRepository(t)
	taskService := NewTaskService(mockTaskRepo)

	// Test data with past due date
	pastDate := time.Now().Add(-24 * time.Hour)
	req := dto.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		DueDate:     &dto.JSONTime{Time: pastDate},
	}

	// Execute
	task, err := taskService.Create(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Equal(t, "due date must be in the future", err.Error())
}

func TestTaskService_GetByID_Success(t *testing.T) {
	// Setup
	mockTaskRepo := mocks.NewMockTaskRepository(t)
	taskService := NewTaskService(mockTaskRepo)

	// Test data
	taskID := primitive.NewObjectID()
	expectedTask := &model.Task{
		ID:          taskID,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      model.TaskStatusPending,
		Priority:    model.TaskPriorityMedium,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Mock expectations
	mockTaskRepo.EXPECT().
		FindByID(mock.Anything, taskID).
		Return(expectedTask, nil).
		Once()

	// Execute
	task, err := taskService.GetByID(context.Background(), taskID.Hex())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, expectedTask.ID, task.ID)
	assert.Equal(t, expectedTask.Title, task.Title)
}

func TestTaskService_GetByID_InvalidID(t *testing.T) {
	// Setup
	mockTaskRepo := mocks.NewMockTaskRepository(t)
	taskService := NewTaskService(mockTaskRepo)

	// Execute with invalid ID
	task, err := taskService.GetByID(context.Background(), "invalid-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Equal(t, "invalid task ID", err.Error())
}

func TestTaskService_GetByID_NotFound(t *testing.T) {
	// Setup
	mockTaskRepo := mocks.NewMockTaskRepository(t)
	taskService := NewTaskService(mockTaskRepo)

	// Test data
	taskID := primitive.NewObjectID()

	// Mock expectations
	mockTaskRepo.EXPECT().
		FindByID(mock.Anything, taskID).
		Return(nil, nil).
		Once()

	// Execute
	task, err := taskService.GetByID(context.Background(), taskID.Hex())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Equal(t, "task not found", err.Error())
}

func TestTaskService_List_Success(t *testing.T) {
	// Setup
	mockTaskRepo := mocks.NewMockTaskRepository(t)
	taskService := NewTaskService(mockTaskRepo)

	// Test data
	params := dto.TaskQueryParams{
		Page:      1,
		Limit:     10,
		Status:    "pending",
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	expectedTasks := []model.Task{
		{
			ID:          primitive.NewObjectID(),
			Title:       "Task 1",
			Description: "Description 1",
			Status:      model.TaskStatusPending,
			Priority:    model.TaskPriorityHigh,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Task 2",
			Description: "Description 2",
			Status:      model.TaskStatusPending,
			Priority:    model.TaskPriorityMedium,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Mock expectations
	mockTaskRepo.EXPECT().
		Find(mock.Anything, mock.MatchedBy(func(filters repository.TaskFilters) bool {
			return filters.Page == 1 &&
				filters.Limit == 10 &&
				filters.Status == "pending"
		})).
		Return(expectedTasks, int64(2), nil).
		Once()

	// Execute
	tasks, meta, err := taskService.List(context.Background(), params)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, int64(2), meta.Total)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 10, meta.Limit)
	assert.Equal(t, 1, meta.TotalPages)
}

func TestTaskService_Update_Success(t *testing.T) {
	// Setup
	mockTaskRepo := mocks.NewMockTaskRepository(t)
	taskService := NewTaskService(mockTaskRepo)

	// Test data
	taskID := primitive.NewObjectID()
	existingTask := &model.Task{
		ID:          taskID,
		Title:       "Old Title",
		Description: "Old Description",
		Status:      model.TaskStatusPending,
		Priority:    model.TaskPriorityLow,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	updateReq := dto.UpdateTaskRequest{
		Title:    "New Title",
		Status:   "completed",
		Priority: "high",
	}

	// Mock expectations
	mockTaskRepo.EXPECT().
		FindByID(mock.Anything, taskID).
		Return(existingTask, nil).
		Once()

	mockTaskRepo.EXPECT().
		Update(mock.Anything, mock.MatchedBy(func(task *model.Task) bool {
			return task.Title == "New Title" &&
				task.Status == model.TaskStatusCompleted &&
				task.Priority == model.TaskPriorityHigh
		})).
		Return(nil).
		Once()

	// Execute
	updatedTask, err := taskService.Update(context.Background(), taskID.Hex(), updateReq)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedTask)
	assert.Equal(t, "New Title", updatedTask.Title)
	assert.Equal(t, model.TaskStatusCompleted, updatedTask.Status)
	assert.Equal(t, model.TaskPriorityHigh, updatedTask.Priority)
}

func TestTaskService_Update_PastDueDate(t *testing.T) {
	// Setup
	mockTaskRepo := mocks.NewMockTaskRepository(t)
	taskService := NewTaskService(mockTaskRepo)

	// Test data
	futureDate := time.Now().Add(24 * time.Hour)
	taskID := primitive.NewObjectID()
	existingTask := &model.Task{
		ID:          taskID,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      model.TaskStatusPending,
		DueDate:     &futureDate,
		Priority:    model.TaskPriorityMedium,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	pastDate := time.Now().Add(-24 * time.Hour)
	updateReq := dto.UpdateTaskRequest{
		DueDate: &dto.JSONTime{Time: pastDate},
	}

	// Mock expectations
	mockTaskRepo.EXPECT().
		FindByID(mock.Anything, taskID).
		Return(existingTask, nil).
		Once()

	// Execute
	updatedTask, err := taskService.Update(context.Background(), taskID.Hex(), updateReq)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updatedTask)
	assert.Equal(t, "due date must be in the future", err.Error())
}

func TestTaskService_Delete_Success(t *testing.T) {
	// Setup
	mockTaskRepo := mocks.NewMockTaskRepository(t)
	taskService := NewTaskService(mockTaskRepo)

	// Test data
	taskID := primitive.NewObjectID()

	// Mock expectations
	mockTaskRepo.EXPECT().
		Delete(mock.Anything, taskID).
		Return(nil).
		Once()

	// Execute
	err := taskService.Delete(context.Background(), taskID.Hex())

	// Assert
	assert.NoError(t, err)
}

func TestTaskService_Delete_InvalidID(t *testing.T) {
	// Setup
	mockTaskRepo := mocks.NewMockTaskRepository(t)
	taskService := NewTaskService(mockTaskRepo)

	// Execute with invalid ID
	err := taskService.Delete(context.Background(), "invalid-id")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid task ID", err.Error())
}

func TestTaskService_Delete_NotFound(t *testing.T) {
	// Setup
	mockTaskRepo := mocks.NewMockTaskRepository(t)
	taskService := NewTaskService(mockTaskRepo)

	// Test data
	taskID := primitive.NewObjectID()

	// Mock expectations
	mockTaskRepo.EXPECT().
		Delete(mock.Anything, taskID).
		Return(errors.New("task not found")).
		Once()

	// Execute
	err := taskService.Delete(context.Background(), taskID.Hex())

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "task not found", err.Error())
}
