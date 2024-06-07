package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ZnNr/todo-list/internal/database"
	taskerror "github.com/ZnNr/todo-list/internal/error"
	"github.com/ZnNr/todo-list/internal/model"
	"github.com/ZnNr/todo-list/internal/task"
	"net/http"
	"strconv"
	"time"
)

var TaskServiceInstance task.TaskService

// taskFromRequestBody извлекает задачу из тела запроса
func taskFromRequestBody(r *http.Request) (model.Task, error) {
	var task model.Task

	buff := bytes.Buffer{}

	_, err := buff.ReadFrom(r.Body)
	if err != nil {
		return model.Task{}, err
	}

	err = json.Unmarshal(buff.Bytes(), &task)
	return task, err
}

func DonePostTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	err := TaskServiceInstance.DoneTask(id)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}
	w.Write([]byte("{}"))
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	err := TaskServiceInstance.DeleteTask(id)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}
	w.Write([]byte("{}"))
}

// PostTask обрабатывает POST запрос для создания задачи
func PostTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	task, err := taskFromRequestBody(r)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}

	id, err := TaskServiceInstance.CreateTask(task)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}

	responseBody, err := json.Marshal(struct {
		Id int `json:"id"`
	}{Id: id})
	if err != nil {
		writeErrorAndRespond(w, http.StatusInternalServerError, err)
		return
	}

	_, err = w.Write(responseBody)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}
}

// GetTask обрабатывает запрос на получение задачи по ID
func GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// Получаем значение параметра "id" из URL запроса
	id := r.URL.Query().Get("id")
	// Получаем задачу по ID с помощью сервиса TaskServiceInstance
	task, err := TaskServiceInstance.GetTask(id)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}
	// Преобразуем полученную задачу в формат JSON и отправляем клиенту
	response, err := json.Marshal(task)
	if err != nil {
		writeErrorAndRespond(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(response)
}

// GetTasks обрабатывает запрос на получение списка задач
func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// Получаем значение параметра "search" из URL запроса
	//search := r.URL.Query().Get("search")
	date := r.URL.Query().Get("date")
	status := r.URL.Query().Get("status")
	pageStr := r.URL.Query().Get("page")
	itemsPerPageStr := r.URL.Query().Get("itemsPerPage")

	page, _ := strconv.Atoi(pageStr)
	itemsPerPage, _ := strconv.Atoi(itemsPerPageStr)
	var tasks *model.TaskList
	var err error
	// Если параметр "search" не указан, получаем все задачи, иначе ищем задачи по запросу
	if len(date) == 0 {
		tasks, err = TaskServiceInstance.GetTasksByStatus(status, page, itemsPerPage)
		if err != nil {
			http.Error(w, "Failed to get tasks by status", http.StatusInternalServerError)
			return
		}
	} else {
		tasks, err = TaskServiceInstance.GetTasksByDateAndStatus(date, status, page, itemsPerPage)
		if err != nil {
			http.Error(w, "Failed to get tasks by date and status", http.StatusInternalServerError)
			return
		}
	}
	if err != nil {
		writeErrorAndRespond(w, http.StatusInternalServerError, err)
		return
	}
	// Преобразуем список задач в формат JSON и отправляем клиенту
	response, err := json.Marshal(tasks)
	if err != nil {
		writeErrorAndRespond(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(response)
}

func PutTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	task, err := taskFromRequestBody(r)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}

	err = TaskServiceInstance.UpdateTask(task)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}
	w.Write([]byte("{}"))
}

// CheckDatabaseAvailability выполняет проверку доступности базы данных (Ping).
func CheckDatabaseAvailability(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	db, err := database.NewPostgresDB()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		http.Error(w, "Database ping error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("pong")); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// writeErrorAndRespond пишет ошибку в ответ и устанавливает соответствующий код состояния
func writeErrorAndRespond(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Write(taskerror.MarshalError(err))
}
