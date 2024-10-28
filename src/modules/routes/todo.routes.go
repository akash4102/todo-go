package routes

import (
	"todo/src/modules/controller"

	"github.com/go-chi/chi/v5"
)

func TodoRoutes(router chi.Router, todoController *controller.TodoController) {
	router.Route("/todos", func(r chi.Router) {
		r.Get("/", todoController.GetAllTodos)
		r.Get("/{id}", todoController.GetTodoByID)
		r.Post("/", todoController.CreateTodo)
		r.Put("/{id}", todoController.UpdateTodo)
		r.Delete("/{id}", todoController.DeleteTodo)
	})
}
