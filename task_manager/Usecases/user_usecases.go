package Usecases

import (
	"errors"

	"task_manager/Domain"
	"task_manager/Infrastructure"
	"task_manager/Repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrBadCreds = errors.New("invalid credentials")

type UserUsecases interface {
	Register(username, password string) (Domain.User, error)
	Login(username, password string) (Domain.User, string, error)
	Promote(username string) (Domain.User, error)
}

type UserUsecasesImpl struct {
	repo        Repositories.UserRepository
	passSvc     Infrastructure.PasswordService
	jwtSvc      Infrastructure.JWTService
}

func NewUserUsecases(repo Repositories.UserRepository, passSvc Infrastructure.PasswordService, jwtSvc Infrastructure.JWTService) UserUsecases {
	return &UserUsecasesImpl{
		repo:    repo,
		passSvc: passSvc,
		jwtSvc:  jwtSvc,
	}
}

func (u *UserUsecasesImpl) Register(username, password string) (Domain.User, error) {
	if username == "" || password == "" {
		return Domain.User{}, errors.New("username and password are required")
	}

	count, err := u.repo.CountUsers()
	if err != nil {
		return Domain.User{}, err
	}

	role := Domain.RoleUser
	if count == 0 {
		role = Domain.RoleAdmin
	}

	hash, err := u.passSvc.Hash(password)
	if err != nil {
		return Domain.User{}, err
	}

	user := Domain.User{
		ID:           primitive.NewObjectID().Hex(),
		Username:     username,
		PasswordHash: hash,
		Role:         role,
	}

	return u.repo.Create(user)
}

func (u *UserUsecasesImpl) Login(username, password string) (Domain.User, string, error) {
	user, err := u.repo.FindByUsername(username)
	if err != nil {
		return Domain.User{}, "", ErrBadCreds
	}

	if err := u.passSvc.Compare(user.PasswordHash, password); err != nil {
		return Domain.User{}, "", ErrBadCreds
	}

	token, err := u.jwtSvc.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		return Domain.User{}, "", err
	}

	return user, token, nil
}

func (u *UserUsecasesImpl) Promote(username string) (Domain.User, error) {
	return u.repo.UpdateRole(username, Domain.RoleAdmin)
}
