package repository

import (
	"context"

	"github.com/grachmannico95/mileapp-test-be/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}
