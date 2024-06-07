package model

type Task struct {
	Id string `json:"id"`

	Title string `json:"title,omitempty"`

	Description string `json:"description,omitempty"`

	Date string `json:"date,omitempty"`

	Status string `json:"status,omitempty"`
}

// TaskList Структура представляет собой список задач
type TaskList struct {
	Tasks []Task `json:"tasks"`
}
