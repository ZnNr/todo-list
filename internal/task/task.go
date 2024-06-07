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

//func (service TaskService) GetTasks() (*model.TaskList, error) {
//	list, err := service.taskData.GetTasks(settings.TasksListRowsLimit)
//	if err != nil {
//		return nil, err
//	}
//	return sliceToTasks(list), err
//}

func (service TaskService) GetTasksByStatus(status string, page int, itemsPerPage int) (*model.TaskList, error) {
	list, err := service.taskData.GetTasksByStatus(status, page, itemsPerPage)
	return sliceToTasks(list), err
}

func (service TaskService) GetTasksByDateAndStatus(date string, status string, page int, itemsPerPage int) (*model.TaskList, error) {
	list, err := service.taskData.GetTasksByDateAndStatus(date, status, page, itemsPerPage)
	return sliceToTasks(list), err
}

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
