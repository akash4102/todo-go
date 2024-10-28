package controller

import (
	"encoding/json"
	"net/http"
	"todo/src/modules/models"
	"todo/src/modules/services"
	"todo/src/pkg/response"

	"github.com/go-chi/chi/v5"
)

type TodoController struct {
	service *services.TodoService
}

func NewTodoController(service *services.TodoService) *TodoController {
	return &TodoController{service: service}
}

func (h *TodoController) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos := h.service.GetAllTodos()
	response.JSON(w, http.StatusOK, todos)
}

func (h *TodoController) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	todo, err := h.service.GetTodoByID(id)
	if err != nil {
		response.JSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	response.JSON(w, http.StatusOK, todo)
}

func (h *TodoController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	createdTodo := h.service.CreateTodo(&todo)
	response.JSON(w, http.StatusCreated, createdTodo)
}

func (h *TodoController) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	todo.ID = id
	if err := h.service.UpdateTodo(id, &todo); err != nil {
		response.JSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	response.JSON(w, http.StatusOK, todo)
}

func (h *TodoController) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.DeleteTodo(id); err != nil {
		response.JSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	response.JSON(w, http.StatusNoContent, nil)
}
