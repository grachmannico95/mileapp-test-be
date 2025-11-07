package repository

import (
	"context"
	"errors"
	"time"

	"github.com/grachmannico95/mileapp-test-be/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type taskRepositoryImpl struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) TaskRepository {
	return &taskRepositoryImpl{
		collection: db.Collection("tasks"),
	}
}

func (r *taskRepositoryImpl) Create(ctx context.Context, task *model.Task) error {
	task.ID = primitive.NewObjectID()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, task)
	return err
}

func (r *taskRepositoryImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Task, error) {
	var task model.Task
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}

func (r *taskRepositoryImpl) Find(ctx context.Context, filters TaskFilters) ([]model.Task, int64, error) {
	query := bson.M{}

	if filters.Status != "" {
		query["status"] = filters.Status
	}

	if filters.Priority != "" {
		query["priority"] = model.PriorityStringToInt(filters.Priority)
	}

	if filters.Search != "" {
		query["$or"] = []bson.M{
			{"title": bson.M{"$regex": filters.Search, "$options": "i"}},
			{"description": bson.M{"$regex": filters.Search, "$options": "i"}},
		}
	}

	if filters.DueDateFrom != "" || filters.DueDateTo != "" {
		dateQuery := bson.M{}

		layout := "2006-01-02"
		if filters.DueDateFrom != "" {
			fromDate, err := time.Parse(layout, filters.DueDateFrom)
			if err == nil {
				dateQuery["$gte"] = fromDate
			}
		}

		if filters.DueDateTo != "" {
			toDate, err := time.Parse(layout, filters.DueDateTo)
			if err == nil {
				dateQuery["$lte"] = toDate
			}
		}

		if len(dateQuery) > 0 {
			query["due_date"] = dateQuery
		}
	}

	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	sortField := "created_at"
	if filters.SortBy != "" {
		sortField = filters.SortBy
	}

	sortOrder := -1
	if filters.SortOrder == "asc" {
		sortOrder = 1
	}

	page := filters.Page
	if page < 1 {
		page = 1
	}

	limit := filters.Limit
	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	findOptions := options.Find().
		SetSort(bson.D{{Key: sortField, Value: sortOrder}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, 0, err
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	return tasks, total, nil
}

func (r *taskRepositoryImpl) Update(ctx context.Context, task *model.Task) error {
	task.UpdatedAt = time.Now()

	filter := bson.M{"_id": task.ID}
	update := bson.M{"$set": task}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (r *taskRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}
