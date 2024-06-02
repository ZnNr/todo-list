package task

import (
	"github.com/ZnNr/todo-list/internal/database"
	"github.com/ZnNr/todo-list/internal/error"
	"github.com/ZnNr/todo-list/internal/model"
	"github.com/ZnNr/todo-list/internal/settings"
	"strconv"
	"time"
)

// TaskService представляет сервис для работы с задачами
type TaskService struct {
	taskData *database.TaskData
}

// sliceToTasks конвертирует срез задач в структуру TaskList
func sliceToTasks(list []model.Task) *model.TaskList {
	if list == nil {
		return &model.TaskList{Tasks: []model.Task{}}
	}
	return &model.TaskList{Tasks: list}
}

// Функция convertTask конвертирует и проверяет задачу перед сохранением
func convertTask(task *model.Task) error {
	if len(task.Title) == 0 {
		return taskerror.ErrRequireTitle
	}
	// Установка даты по умолчанию, если она не была указана, и проверка формата даты
	now := time.Now().Format(settings.DateFormat)
	if len(task.Date) == 0 {
		task.Date = now
	}
	_, err := time.Parse(settings.DateFormat, task.Date)
	if err != nil {
		return err
	}
	return nil
}

// InitTaskService создает новый экземпляр TaskService
func InitTaskService(taskData *database.TaskData) TaskService {
	return TaskService{taskData: taskData}
}

// CreateTask Метод создает новую задачу
func (service TaskService) CreateTask(task model.Task) (int, error) {
	err := convertTask(&task)
	if err != nil {
		return 0, err
	}
	id, err := service.taskData.InsertTask(task)
	return int(id), err
}

// UpdateTask обновляет существующую задач
func (service TaskService) UpdateTask(task model.Task) error {
	err := convertTask(&task)
	if err != nil {
		return err
	}

	updated, err := service.taskData.UpdateTask(task)
	if err != nil {
		return err
	}
	if !updated {
		return taskerror.ErrNotFoundTask
	}
	return nil
}

// GetTasks возвращает список задач
func (service TaskService) GetTasks() (*model.TaskList, error) {
	list, err := service.taskData.GetTasks(settings.TasksListRowsLimit)
	if err != nil {
		return nil, err
	}
	return sliceToTasks(list), err
}

// GetTasksByStatus возвращает задачи по статусу
func (service TaskService) GetTasksByStatus(status string) (*model.TaskList, error) {
	list, err := service.taskData.GetTasksByStatus(status, settings.TasksListRowsLimit)
	return sliceToTasks(list), err
}

// GetTasksByDateAndStatus возвращает задачи по дате и статусу
func (service TaskService) GetTasksByDateAndStatus(status string, date string) (*model.TaskList, error) {
	list, err := service.taskData.GetTasksByDateAndStatus(status, date, settings.TasksListRowsLimit)
	return sliceToTasks(list), err
}

// GetTask возвращает задачу по идентификатору
func (service TaskService) GetTask(id string) (*model.Task, error) {
	convId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	task, err := service.taskData.GetTask(convId)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// DeleteTask удаляет задачу по идентификатору
func (service TaskService) DeleteTask(id string) error {
	convId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	deleted, err := service.taskData.DeleteTask(convId)
	if err != nil {
		return err
	}
	if !deleted {
		return taskerror.ErrNotFoundTask
	}
	return nil
}

// DoneTask помечает задачу как выполненную
func (service TaskService) DoneTask(id string) error {
	convId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	task, err := service.taskData.GetTask(convId)
	if err != nil {
		return err
	}

	updated, err := service.taskData.UpdateTask(task)
	if err != nil {
		return err
	}
	if !updated {
		return taskerror.ErrNotFoundTask
	}
	return nil
}
