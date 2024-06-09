package router

import (
	"github.com/ZnNr/todo-list/internal/handlers"
	"github.com/ZnNr/todo-list/internal/settings"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func StartServer() {

	// Инициализация маршрутизатора.
	r := chi.NewRouter()
	// Установка маршрутов для обработки файлов и API.

	// Группировка маршрутов API для задач.
	r.Group(func(r chi.Router) {
		r.Post("/task", handlers.PostTask)                     // Создание задачи
		r.Put("/task", handlers.PutTask)                       // Обновление задачи
		r.Delete("/task", handlers.DeleteTask)                 // Удаление задачи
		r.Get("/task", handlers.GetTask)                       // Получение конкретной задачи
		r.Post("/task/done", handlers.DonePostTask)            // Отметка задачи как выполненной
		r.Get("/tasks", handlers.GetALLTasks)                  // API для получения списка ВСЕХ задач
		r.Get("/status", handlers.GetTasksByStatus)            // API для получения списка списка всех задачь с определенныфм статусом
		r.Get("/datestatus", handlers.GetTasksByDateAndStatus) // API для получения списка ВСЕХ задач с определенным статусом и за определенную дату

	})

	// Старт веб-сервера на указанном порту.
	port := settings.Setting("TODO_PORT")
	serverAddr := ":" + port
	log.Printf("Starting server on %s...", serverAddr)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Error starting the web server: %v", err)
	}
}
