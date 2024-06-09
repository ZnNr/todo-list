package main

import (
	"github.com/ZnNr/todo-list/internal/database"
	"github.com/ZnNr/todo-list/internal/handlers"
	"github.com/ZnNr/todo-list/internal/router"
	"github.com/ZnNr/todo-list/internal/settings"
	"github.com/ZnNr/todo-list/internal/task"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	// Инициализация базы данных и задач.
	dbFile := settings.Setting("TODO_DBFILE")
	taskData, dbErr := database.NewTaskData(dbFile)
	defer taskData.CloseDb()
	if dbErr != nil {
		log.Fatalf("Error initializing task data: %v", dbErr)
	}

	// Инициализация службы задач.
	handlers.TaskServiceInstance = task.InitTaskService(taskData)

	// Инициализация маршрутизатора и запуск сервера.
	router.StartServer()
}
