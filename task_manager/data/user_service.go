package data

import (
	"context"
	"errors"
	"time"

	"task_manager/models"

	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrBadCreds     = errors.New("invalid credentials")
	ErrUsersCollectionEmptyCheck = errors.New("failed to check users count")
)

type UserService struct {
	collection *mongo.Collection
}

func NewUserService(col *mongo.Collection) *UserService {
	return &UserService{collection: col}
}

func (s *UserService) CountUsers() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.collection.CountDocuments(ctx, bson.M{})
}

func (s *UserService) FindByUsername(username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.User{}, ErrNotFound
	}
	return user, err
}

func (s *UserService) CreateUser(username, password, role string) (models.User, error) {
	// Ensure username unique
	if _, err := s.FindByUsername(username); err == nil {
		return models.User{}, ErrUserExists
	} else if err != nil && err != ErrNotFound {
		return models.User{}, err
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		ID:           primitive.NewObjectID().Hex(),
		Username:     username,
		PasswordHash: string(hashBytes),
		Role:         role,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = s.collection.InsertOne(ctx, user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserService) Authenticate(username, password string) (models.User, error) {
	user, err := s.FindByUsername(username)
	if err != nil {
		return models.User{}, ErrBadCreds
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return models.User{}, ErrBadCreds
	}
	return user, nil
}

func (s *UserService) PromoteToAdmin(username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update role
	res := s.collection.FindOneAndUpdate(
		ctx,
		bson.M{"username": username},
		bson.M{"$set": bson.M{"role": "admin"}},
	)
	var updated models.User
	err := res.Decode(&updated)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.User{}, ErrNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	updated.Role = "admin"
	return updated, nil
}
