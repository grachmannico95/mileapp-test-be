package repository

import (
	"context"

	"github.com/grachmannico95/mileapp-test-be/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskFilters struct {
	Status      string
	Priority    string
	Search      string
	DueDateFrom string
	DueDateTo   string
	SortBy      string
	SortOrder   string
	Page        int
	Limit       int
}

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Task, error)
	Find(ctx context.Context, filters TaskFilters) ([]model.Task, int64, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
