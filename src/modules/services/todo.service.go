package services

import (
	"context"
	"todo/src/modules/models"
	"todo/src/modules/repository"

	"go.mongodb.org/mongo-driver/bson"
)

type TodoService struct {
	repo *repository.MongoRepo
}

func NewTodoService(repo *repository.MongoRepo) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) GetAllTodos() ([]*models.Todo, error) {
	return s.repo.GetAll()
}

func (s *TodoService) GetTodoByID(id string) (*models.Todo, error) {
	return s.repo.GetByID(id)
}

func (s *TodoService) CreateTodo(todo *models.Todo) (*models.Todo, error) {
	return s.repo.Create(todo)
}

func (s *TodoService) UpdateTodo(id string, todo *models.Todo) error {
	return s.repo.Update(id, todo)
}

func (s *TodoService) DeleteTodo(id string) error {
	return s.repo.Delete(id)
}

func (s *TodoService) GetTodoMetrics(ctx context.Context) ([]bson.M, error) {
	return s.repo.GetTodoMetrics(ctx)
}
