package data

import (
	"errors"
	"sync"

	"task_manager/models"
)

type TaskService struct {
	mu     sync.RWMutex
	tasks  map[int64]models.Task
	nextID int64
}

// creates a new in memory task serrvice
func NewTaskService() *TaskService {
	return &TaskService{
		tasks:  make(map[int64]models.Task),
		nextID: 1,
	}
}

// returns all tasks
func (s *TaskService) GetAll() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}
	return result
}

// returns a task by ID
func (s *TaskService) GetByID(id int64) (models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, errors.New("task not found")
	}
	return task, nil
}

// adds a new task and returns assigned id
func (s *TaskService) Create(task models.Task) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = s.nextID
	s.nextID++
	s.tasks[task.ID] = task
	return task
}

// updates an existing task
func (s *TaskService) Update(id int64, updated models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.tasks[id]
	if !ok {
		return models.Task{}, errors.New("task not found")
	}

	updated.ID = id // preserve ID
	s.tasks[id] = updated
	return updated, nil
}

// removes a task by ID
func (s *TaskService) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return errors.New("task not found")
	}
	delete(s.tasks, id)
	return nil
}
