package Repositories

import (
	"context"
	"errors"
	"time"

	"task_manager/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrUserExists = errors.New("user already exists")

type UserRepository interface {
	CountUsers() (int64, error)
	FindByUsername(username string) (Domain.User, error)
	Create(user Domain.User) (Domain.User, error)
	UpdateRole(username string, role Domain.Role) (Domain.User, error)
}

type MongoUserRepository struct {
	col *mongo.Collection
}

func NewMongoUserRepository(col *mongo.Collection) UserRepository {
	return &MongoUserRepository{col: col}
}

func (r *MongoUserRepository) CountUsers() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.col.CountDocuments(ctx, bson.M{})
}

func (r *MongoUserRepository) FindByUsername(username string) (Domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user Domain.User
	err := r.col.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Domain.User{}, ErrNotFound
	}
	return user, err
}

func (r *MongoUserRepository) Create(user Domain.User) (Domain.User, error) {
	// Ensure unique username
	if _, err := r.FindByUsername(user.Username); err == nil {
		return Domain.User{}, ErrUserExists
	} else if err != nil && err != ErrNotFound {
		return Domain.User{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.col.InsertOne(ctx, user)
	if err != nil {
		return Domain.User{}, err
	}
	return user, nil
}

func (r *MongoUserRepository) UpdateRole(username string, role Domain.Role) (Domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res := r.col.FindOneAndUpdate(
		ctx,
		bson.M{"username": username},
		bson.M{"$set": bson.M{"role": string(role)}},
	)

	var updated Domain.User
	err := res.Decode(&updated)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Domain.User{}, ErrNotFound
	}
	if err != nil {
		return Domain.User{}, err
	}

	updated.Role = role
	return updated, nil
}
