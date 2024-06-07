package database

import (
	"database/sql"
	"fmt"
	"github.com/ZnNr/todo-list/internal/model"
	_ "github.com/lib/pq"
)

const (
	driverName = "postgres"

	tableSchema = `
 CREATE TABLE IF NOT EXISTS todolist ( 
     id SERIAL PRIMARY KEY, 
     date DATE, 
     title TEXT, 
     description TEXT, 
     status TEXT );
`
	indexSchema = `
CREATE INDEX IF NOT EXISTS indexdate ON todolist (date);
`
	insertQuery = `
INSERT INTO todolist(date, title, description) VALUES ($1, $2, $3)
`
	getTaskQuery = "SELECT * FROM todolist WHERE id = $1"

	getTasksQuery = "SELECT * FROM todolist ORDER BY date LIMIT $1"

	updateQuery = "UPDATE todolist SET date=$1, title=$2, description=$3 WHERE id=$4"

	deleteQuery = "DELETE FROM todolist WHERE id=$1"
)

// TaskData представляет структуру для работы с данными задач
type TaskData struct {
	db *sql.DB
}

// NewTaskData создает новый экземпляр TaskData с подключением к базе данных
func NewTaskData(db *sql.DB) *TaskData {
	return &TaskData{db: db}
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

func (data TaskData) GetTasksByDateAndStatus(date string, status string, page int, itemsPerPage int) ([]model.Task, error) {
	offset := (page - 1) * itemsPerPage

	// Генерация SQL запроса с учетом фильтрации по дате и статусу, а также пагинации
	query := "SELECT * FROM todolist WHERE 1=1"
	var params []interface{}
	if date != "" {
		query += " AND date = $1"
		params = append(params, date)
	}
	if status == "Не выполнено" || status == "Выполнено" {
		query += " AND status = $2"
		params = append(params, status)
	}
	query += fmt.Sprintf(" ORDER BY date LIMIT %d OFFSET %d", itemsPerPage, offset)

	// Выполнение запроса к базе данных
	rows, err := data.db.Query(query, params...)
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

func (data TaskData) GetTasksByStatus(status string, page int, itemsPerPage int) ([]model.Task, error) {
	offset := (page - 1) * itemsPerPage

	// Генерация SQL запроса с учетом фильтрации по дате и статусу, а также пагинации
	query := "SELECT * FROM todolist WHERE 1=1"
	var params []interface{}

	if status == "Не выполнено" || status == "Выполнено" {
		query += " AND status = $2"
		params = append(params, status)
	}
	query += fmt.Sprintf(" ORDER BY date LIMIT %d OFFSET %d", itemsPerPage, offset)

	// Выполнение запроса к базе данных
	rows, err := data.db.Query(query, params...)
	if err != nil {
		return nil, err
	}

	return getTasksByRows(rows)
}

// openDb открывает соединение с базой данных
func openDb(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if _, err := tx.Exec(tableSchema); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(indexSchema); err != nil {
		return nil, err
	}

	return db, nil
}
