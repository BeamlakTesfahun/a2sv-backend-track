package Repositories

import (
	"context"
	"errors"
	"time"

	"task_manager/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrNotFound = errors.New("not found")

type TaskRepository interface {
	GetAll() ([]Domain.Task, error)
	GetByID(id string) (Domain.Task, error)
	Create(task Domain.Task) (Domain.Task, error)
	Update(id string, task Domain.Task) (Domain.Task, error)
	Delete(id string) error
}

type MongoTaskRepository struct {
	col *mongo.Collection
}

func NewMongoTaskRepository(col *mongo.Collection) TaskRepository {
	return &MongoTaskRepository{col: col}
}

func (r *MongoTaskRepository) GetAll() ([]Domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tasks []Domain.Task
	if err := cur.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *MongoTaskRepository) GetByID(id string) (Domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var task Domain.Task
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Domain.Task{}, ErrNotFound
	}
	return task, err
}

func (r *MongoTaskRepository) Create(task Domain.Task) (Domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.col.InsertOne(ctx, task)
	if err != nil {
		return Domain.Task{}, err
	}
	return task, nil
}

func (r *MongoTaskRepository) Update(id string, task Domain.Task) (Domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	task.ID = id
	res, err := r.col.ReplaceOne(ctx, bson.M{"_id": id}, task)
	if err != nil {
		return Domain.Task{}, err
	}
	if res.MatchedCount == 0 {
		return Domain.Task{}, ErrNotFound
	}
	return task, nil
}

func (r *MongoTaskRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
