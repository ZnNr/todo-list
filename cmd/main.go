package main

import (
	"github.com/ZnNr/todo-list/internal/database"
	"github.com/ZnNr/todo-list/internal/handlers"
	"github.com/ZnNr/todo-list/internal/router"
	"github.com/ZnNr/todo-list/internal/task"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	// Инициализация базы данных и задач.
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	taskData := database.NewTaskData(db)

	// Инициализация службы задач.
	handlers.TaskServiceInstance = task.InitTaskService(taskData)

	// Инициализация маршрутизатора и запуск сервера.
	router.StartServer()
}
