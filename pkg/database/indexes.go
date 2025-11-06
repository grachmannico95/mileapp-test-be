package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupIndexes(ctx context.Context, db *mongo.Database) error {
	usersCollection := db.Collection("users")

	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	if _, err := usersCollection.Indexes().CreateOne(ctx, emailIndex); err != nil {
		return fmt.Errorf("failed to create email index: %w", err)
	}

	tasksCollection := db.Collection("tasks")

	statusIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}},
	}

	if _, err := tasksCollection.Indexes().CreateOne(ctx, statusIndex); err != nil {
		return fmt.Errorf("failed to create status index: %w", err)
	}

	priorityIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "priority", Value: 1}},
	}

	if _, err := tasksCollection.Indexes().CreateOne(ctx, priorityIndex); err != nil {
		return fmt.Errorf("failed to create priority index: %w", err)
	}

	dueDateIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "due_date", Value: 1}},
	}

	if _, err := tasksCollection.Indexes().CreateOne(ctx, dueDateIndex); err != nil {
		return fmt.Errorf("failed to create due_date index: %w", err)
	}

	createdAtIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: -1}},
	}

	if _, err := tasksCollection.Indexes().CreateOne(ctx, createdAtIndex); err != nil {
		return fmt.Errorf("failed to create created_at index: %w", err)
	}

	return nil
}
