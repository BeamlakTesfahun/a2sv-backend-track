package data

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"task_manager/models"
)

var (
	// ErrNotFound is returned when a task does not exist.
	ErrNotFound = errors.New("task not found")
)

// TaskService handles MongoDB-backed task operations.
type TaskService struct {
	collection *mongo.Collection
}

// NewTaskService creates a new TaskService using the given MongoDB collection.
func NewTaskService(collection *mongo.Collection) *TaskService {
	return &TaskService{
		collection: collection,
	}
}

// contextWithTimeout creates a context with a reasonable timeout for DB ops.
func contextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// GetAll returns all tasks from MongoDB.
func (s *TaskService) GetAll() ([]models.Task, error) {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetByID returns a task by ID.
func (s *TaskService) GetByID(id string) (models.Task, error) {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	var task models.Task
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.Task{}, ErrNotFound
	}
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

// Create inserts a new task into MongoDB and returns it.
func (s *TaskService) Create(task models.Task) (models.Task, error) {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	// If ID is empty, generate a new unique string ID.
	if task.ID == "" {
		task.ID = primitive.NewObjectID().Hex()
	}

	_, err := s.collection.InsertOne(ctx, task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

// Update replaces an existing task document.
func (s *TaskService) Update(id string, updated models.Task) (models.Task, error) {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	updated.ID = id // ensure ID matches path parameter

	res, err := s.collection.ReplaceOne(ctx, bson.M{"_id": id}, updated)
	if err != nil {
		return models.Task{}, err
	}
	if res.MatchedCount == 0 {
		return models.Task{}, ErrNotFound
	}

	return updated, nil
}

// Delete removes a task by ID.
func (s *TaskService) Delete(id string) error {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	res, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}
