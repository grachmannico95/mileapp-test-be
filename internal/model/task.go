package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
)

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Status      TaskStatus         `bson:"status" json:"status"`
	Priority    int                `bson:"priority" json:"priority"`
	DueDate     *time.Time         `bson:"due_date,omitempty" json:"due_date,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

func NewTask(title, description string, status TaskStatus, priority TaskPriority, dueDate *time.Time) *Task {
	now := time.Now()
	return &Task{
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    PriorityStringToInt(string(priority)),
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func IsValidStatus(status string) bool {
	switch TaskStatus(status) {
	case TaskStatusPending, TaskStatusInProgress, TaskStatusCompleted:
		return true
	}
	return false
}

func IsValidPriority(priority string) bool {
	switch TaskPriority(priority) {
	case TaskPriorityLow, TaskPriorityMedium, TaskPriorityHigh:
		return true
	}
	return false
}

func PriorityStringToInt(priority string) int {
	switch TaskPriority(priority) {
	case TaskPriorityLow:
		return 1
	case TaskPriorityMedium:
		return 2
	case TaskPriorityHigh:
		return 3
	default:
		return 0
	}
}

func PriorityIntToString(priority int) string {
	switch priority {
	case 1:
		return string(TaskPriorityLow)
	case 2:
		return string(TaskPriorityMedium)
	case 3:
		return string(TaskPriorityHigh)
	default:
		return ""
	}
}
