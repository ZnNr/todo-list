package database

import (
	"database/sql"

	"github.com/ZnNr/todo-list/internal/model"
	_ "modernc.org/sqlite"
)

const (
	driverName = "sqlite"

	tableSchema = `
CREATE TABLE IF NOT EXISTS todolist (
    id INTEGER PRIMARY KEY,
    date VARCHAR(8),
    title TEXT,
    description TEXT, 
    status TEXT
);
`
	indexSchema = `
CREATE INDEX IF NOT EXISTS indexdate ON todolist (date);
`
	insertQuery = `
INSERT INTO scheduler(date, title, description) VALUES (?, ?, ?)
`
	getTaskQuery = "SELECT * FROM todolist WHERE id = ?"

	getTasksQuery = "SELECT * FROM todolist ORDER BY date LIMIT ?"

	getTasksByDateQuery          = "SELECT * FROM todolist WHERE date = ? ORDER BY date LIMIT ?"
	getTasksByStatusQuery        = "SELECT * FROM todolist WHERE status = ? LIMIT ?"
	getTasksByStatusAndDateQuery = "SELECT * FROM todolist WHERE status = ? AND date = ? ORDER BY date LIMIT ?"

	updateQuery = "UPDATE todolist SET date=?, title=?, description=?, WHERE id=?"

	deleteQuery = "DELETE FROM todolist WHERE id=:id"
)

// TaskData представляет структуру для работы с данными задач
type TaskData struct {
	db *sql.DB
}

// NewTaskData создает новый экземпляр TaskData с подключением к базе данных
func NewTaskData(dataSourceName string) (*TaskData, error) {
	db, err := openDb(dataSourceName)
	if err != nil {
		return nil, err
	}
	return &TaskData{db: db}, nil
}

// CloseDb закрывает соединение с базой данных
func (data *TaskData) CloseDb() error {
	return data.db.Close()
}

// InsertTask вставляет задачу в базу данных и возвращает ее ID
func (data *TaskData) InsertTask(task model.Task) (int64, error) {
	stmt, err := data.db.Prepare(insertQuery)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(task.Date, task.Title, task.Description)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

// getTasksByRows извлекает задачи из результата sql.Rows
func getTasksByRows(rows *sql.Rows) ([]model.Task, error) {
	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Description); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTask получает задачу по ID
func (data TaskData) GetTask(id int) (model.Task, error) {

	row := data.db.QueryRow(getTaskQuery, id)

	var task model.Task
	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Description)
	return task, err
}

// GetTasks получает все задачи с ограничением по количеству
func (data TaskData) GetTasks(limit int) ([]model.Task, error) {

	rows, err := data.db.Query(getTasksQuery, limit)
	if err != nil {
		return nil, err
	}
	return getTasksByRows(rows)
}

// GetTasksByDate получает задачи по дате с ограничением по количеству
func (data TaskData) GetTasksByDate(date string, limit int) ([]model.Task, error) {

	rows, err := data.db.Query(getTasksByDateQuery, date, limit)
	if err != nil {
		return nil, err
	}
	return getTasksByRows(rows)
}

// GetTasksByStatus получает задачи по статусу с ограничением по количеству
func (data TaskData) GetTasksByStatus(status string, limit int) ([]model.Task, error) {

	rows, err := data.db.Query(getTasksByStatusQuery, status, limit)
	if err != nil {
		return nil, err
	}
	return getTasksByRows(rows)
}

// GetTasksByDateAndStatus получает задачи по статусу и дате с ограничением по количеству
func (data TaskData) GetTasksByDateAndStatus(status string, date string, limit int) ([]model.Task, error) {
	rows, err := data.db.Query(getTasksByStatusAndDateQuery, status, date, limit)
	if err != nil {
		return nil, err
	}
	return getTasksByRows(rows)
}

// UpdateTask обновляет задачу в базе данных.
func (data TaskData) UpdateTask(task model.Task) (bool, error) {

	// Начало транзакции.
	tx, err := data.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback() // Откат транзакции в случае ошибки.

	// Выполнение подготовленного запроса внутри транзакции.
	result, err := tx.Exec(updateQuery, task.Date, task.Title, task.Description, task.Id)
	if err != nil {
		return false, err
	}

	// Получение количества обновленных строк.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	// Коммит транзакции, если все операции без ошибок.
	if err = tx.Commit(); err != nil {
		return false, err
	}

	// Проверка, что была обновлена одна строка.
	return rowsAffected == 1, nil
}

// DeleteTask удаляет задачу из базы данных по ID
func (data TaskData) DeleteTask(id int) (bool, error) {
	// Получаем задачу по ID для проверки существования
	_, err := data.GetTask(id)
	if err != nil {
		return false, err
	}

	res, err := data.db.Exec(deleteQuery, sql.Named("id", id))
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted == 1, err
}

// openDb открывает соединение с базой данных
func openDb(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(tableSchema); err != nil {
		return nil, err
	}
	if _, err := db.Exec(indexSchema); err != nil {
		return nil, err
	}
	return db, nil
}
