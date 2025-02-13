package tests

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
	"os"
	"testing"
	"time"
)

var Port = 7540
var DBFile = "../todolist.db"

type Task struct {
	Id          int64  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Date        string `db:"date"`
	Status      string `db:"status"`
}

func count(db *sqlx.DB) (int, error) {
	var count int
	return count, db.Get(&count, `SELECT count(id) FROM todolist`)
}

func openDB(t *testing.T) *sqlx.DB {
	dbfile := DBFile
	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbfile = envFile
	}
	db, err := sqlx.Connect("sqlite", dbfile)
	assert.NoError(t, err)
	return db
}

func TestDB(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	before, err := count(db)
	assert.NoError(t, err)

	today := time.Now().Format(`20060102`)

	res, err := db.Exec(`INSERT INTO todolist (date, title, description) 
	VALUES (?, 'Todo', 'Описание')`, today)
	assert.NoError(t, err)

	id, err := res.LastInsertId()

	var task Task
	err = db.Get(&task, `SELECT * FROM todolist WHERE id=?`, id)

	assert.Equal(t, id, task.Id)
	assert.Equal(t, `Todo`, task.Title)
	assert.Equal(t, `Описание`, task.Description)

	_, err = db.Exec(`DELETE FROM todolist WHERE id = ?`, id)
	assert.NoError(t, err)

	after, err := count(db)
	assert.NoError(t, err)

	assert.Equal(t, before, after)
}

func TestDB2(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	before, err := count(db)
	assert.NoError(t, err)

	today := time.Now().Format(`20060102`)

	res, err := db.Exec(`INSERT INTO todolist (date, title, description) 
	VALUES (?, 'Todo', 'Описание')`, today)
	assert.NoError(t, err)

	id, err := res.LastInsertId()

	var task Task
	err = db.Get(&task, `SELECT * FROM todolist WHERE id=?`, id)

	assert.Equal(t, id, task.Id)
	assert.Equal(t, `Todo`, task.Title)
	assert.Equal(t, `Описание`, task.Description)

	_, err = db.Exec(`DELETE FROM todolist WHERE id = ?`, id)
	assert.NoError(t, err)

	after, err := count(db)
	assert.NoError(t, err)

	assert.Equal(t, before, after)
}
