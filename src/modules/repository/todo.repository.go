package repository

import (
	"errors"
	"sync"
	"todo/src/modules/models"

	"github.com/google/uuid"
)

type InMemoryRepo struct {
	mu    sync.Mutex
	todos map[string]*models.Todo
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		todos: make(map[string]*models.Todo),
	}
}

func (repo *InMemoryRepo) GetAll() []*models.Todo {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var todos []*models.Todo
	for _, todo := range repo.todos {
		todos = append(todos, todo)
	}
	return todos
}

func (repo *InMemoryRepo) GetByID(id string) (*models.Todo, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	todo, exists := repo.todos[id]
	if !exists {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

func (repo *InMemoryRepo) Create(todo *models.Todo) *models.Todo {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	todo.ID = uuid.New().String()
	repo.todos[todo.ID] = todo
	return todo
}

func (repo *InMemoryRepo) Update(id string, updatedTodo *models.Todo) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.todos[id]; !exists {
		return errors.New("todo not found")
	}

	repo.todos[id] = updatedTodo
	return nil
}

func (repo *InMemoryRepo) Delete(id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.todos[id]; !exists {
		return errors.New("todo not found")
	}

	delete(repo.todos, id)
	return nil
}
