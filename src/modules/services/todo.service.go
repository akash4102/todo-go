package services

import (
	"context"
	"todo/src/client/clickhouse"
	"todo/src/modules/models"
	"todo/src/modules/repository"

	"go.mongodb.org/mongo-driver/bson"
)

type TodoService struct {
	repo       *repository.MongoRepo
	clickHouse *clickhouse.ClickHouseRepo
}

func NewTodoService(repo *repository.MongoRepo, clickHouse *clickhouse.ClickHouseRepo) *TodoService {
	return &TodoService{repo: repo, clickHouse: clickHouse}
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

func (s *TodoService) GetTodoMetricsMongodb(ctx context.Context) ([]bson.M, error) {
	return s.repo.GetTodoMetrics(ctx)
}

func (s *TodoService) GetTodoMetricsClickHouse(ctx context.Context) ([]map[string]interface{}, error) {
	return s.clickHouse.GetTodoMetrics(ctx)
}
