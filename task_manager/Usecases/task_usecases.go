package Usecases

import (
	"errors"

	"task_manager/Domain"
	"task_manager/Repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskUsecases interface {
	GetAllTasks() ([]Domain.Task, error)
	GetTaskByID(id string) (Domain.Task, error)
	CreateTask(task Domain.Task) (Domain.Task, error)
	UpdateTask(id string, task Domain.Task) (Domain.Task, error)
	DeleteTask(id string) error
}

type TaskUsecasesImpl struct {
	repo Repositories.TaskRepository
}

func NewTaskUsecases(repo Repositories.TaskRepository) TaskUsecases {
	return &TaskUsecasesImpl{repo: repo}
}

func (u *TaskUsecasesImpl) GetAllTasks() ([]Domain.Task, error) {
	return u.repo.GetAll()
}

func (u *TaskUsecasesImpl) GetTaskByID(id string) (Domain.Task, error) {
	return u.repo.GetByID(id)
}

func (u *TaskUsecasesImpl) CreateTask(task Domain.Task) (Domain.Task, error) {
	if task.Title == "" {
		return Domain.Task{}, errors.New("title is required")
	}
	if task.ID == "" {
		task.ID = primitive.NewObjectID().Hex()
	}
	return u.repo.Create(task)
}

func (u *TaskUsecasesImpl) UpdateTask(id string, task Domain.Task) (Domain.Task, error) {
	if task.Title == "" {
		return Domain.Task{}, errors.New("title is required")
	}
	return u.repo.Update(id, task)
}

func (u *TaskUsecasesImpl) DeleteTask(id string) error {
	return u.repo.Delete(id)
}
