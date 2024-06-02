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
	//
	//// Установка маршрутов для обработки файлов и API.
	//r.Get("/*", FileServer) // Обработка запросов к файлам

	// Группировка маршрутов API для задач.
	r.Group(func(r chi.Router) {
		r.Post("/task", handlers.PostTask)          // Создание задачи
		r.Put("/task", handlers.PutTask)            // Обновление задачи
		r.Delete("/task", handlers.DeleteTask)      // Удаление задачи
		r.Get("/task", handlers.GetTask)            // Получение конкретной задачи
		r.Post("/task/done", handlers.DonePostTask) // Отметка задачи как выполненной
		r.Get("/tasks", handlers.GetTasks)          // API для получения списка задач
		r.Get("/ping", handlers.CheckDatabaseAvailability)
	})

	// Старт веб-сервера на указанном порту.
	port := settings.Setting("TODO_PORT")
	serverAddr := ":" + port
	log.Printf("Starting server on %s...", serverAddr)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Error starting the web server: %v", err)
	}
}

//
//// FileServer обрабатывает запросы на статические файлы и отправляет их клиенту.
//func FileServer(w http.ResponseWriter, r *http.Request) {
//	handler := http.FileServer(http.Dir("/" + settings.WebPath)) // Требуется начальный "/"
//	handler.ServeHTTP(w, r)
//}
