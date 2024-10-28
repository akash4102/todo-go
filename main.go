package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todo/src/modules/controller"
	"todo/src/modules/repository"
	"todo/src/modules/routes"
	"todo/src/modules/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Setup in-memory repository and services
	repo := repository.NewInMemoryRepo()

	// Setup services
	todoService := services.NewTodoService(repo)

	// Setup controller
	todoController := controller.NewTodoController(todoService)

	// Setup routes
	routes.TodoRoutes(router, todoController)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error Port is not available in .env file")
		return
	}
	fmt.Printf("Starting server on port %s\n", port)
	http.ListenAndServe(port, router)
}
