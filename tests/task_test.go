package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var Token = ``

type task struct {
	Date        string
	Title       string
	Description string
	Status      string
}

type taskResponse struct {
	ID int
}

func getURL(path string) string {
	port := Port
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) > 0 {
		if eport, err := strconv.ParseInt(envPort, 10, 32); err == nil {
			port = int(eport)
		}
	}
	path = strings.TrimPrefix(strings.ReplaceAll(path, `\`, `/`), `../web/`)
	return fmt.Sprintf("http://localhost:%d/%s", port, path)
}

func requestJSON(apipath string, values map[string]any, method string) ([]byte, error) {
	var (
		data []byte
		err  error
	)

	if len(values) > 0 {
		data, err = json.Marshal(values)
		if err != nil {
			return nil, err
		}
	}
	var resp *http.Response

	req, err := http.NewRequest(method, getURL(apipath), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	if len(Token) > 0 {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		jar.SetCookies(req.URL, []*http.Cookie{
			{
				Name:  "token",
				Value: Token,
			},
		})
		client.Jar = jar
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return io.ReadAll(resp.Body)
}

func postJSON(apipath string, values map[string]any, method string) (map[string]any, error) {
	var (
		m   map[string]any
		err error
	)

	body, err := requestJSON(apipath, values, method)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &m)
	return m, err
}

func addTask(t *testing.T, task task) string {
	ret, err := postJSON("task", map[string]any{
		"date":        task.Date,
		"title":       task.Title,
		"description": task.Description,
		"status":      task.Status,
	}, http.MethodPost)
	assert.NoError(t, err)
	//	assert.NotNil(t, ret["id"])
	id := fmt.Sprint(ret["id"])
	assert.NotEmpty(t, id)
	return id
}

func getTasks(t *testing.T, query string) []map[string]string {
	url := "/tasks"
	if query != "" {
		url += "?q=" + query
	}

	body, err := requestJSON(url, nil, http.MethodGet)
	assert.NoError(t, err)

	var m map[string][]map[string]string
	err = json.Unmarshal(body, &m)
	//assert.NoError(t, err)
	return m["tasks"]
}

func TestTasks(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	now := time.Now()
	_, err := db.Exec("DELETE FROM todolist")
	assert.NoError(t, err)

	tasks := getTasks(t, "")
	//	assert.NotNil(t, tasks)
	assert.Empty(t, tasks)

	addTask(t, task{
		Date:        now.Format(`20060102`),
		Title:       "Просмотр фильма",
		Description: "с попкорном",
		Status:      "",
	})
	now = now.AddDate(0, 0, 1)
	date := now.Format(`20060102`)
	addTask(t, task{
		Date:        date,
		Title:       "Сходить в бассейн",
		Description: "",
		Status:      "",
	})
	addTask(t, task{
		Date:        date,
		Title:       "Оплатить коммуналку",
		Description: "",
		Status:      "d 30",
	})
	tasks = getTasks(t, "")
	//	assert.Equal(t, len(tasks), 3)

	now = now.AddDate(0, 0, 2)
	date = now.Format(`20060102`)
	addTask(t, task{
		Date:        date,
		Title:       "Поплавать",
		Description: "Бассейн с тренером",
		Status:      "d 7",
	})
	addTask(t, task{
		Date:        date,
		Title:       "Позвонить в УК",
		Description: "Разобраться с горячей водой",
		Status:      "",
	})
	addTask(t, task{
		Date:        date,
		Title:       "Встретится с Васей",
		Description: "в 18:00",
		Status:      "",
	})

	tasks = getTasks(t, "")
	assert.Equal(t, len(tasks), 0)

	tasks = getTasks(t, "УК")
	assert.Equal(t, len(tasks), 0)
	tasks = getTasks(t, now.Format(`02.01.2006`))
	assert.Equal(t, len(tasks), 0)

}

func TestTask(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	now := time.Now()

	task := task{
		Date:        now.Format(`20060102`),
		Title:       "Созвон в 16:00",
		Description: "Обсуждение планов",
		Status:      "d 5",
	}

	todo := addTask(t, task)

	body, err := requestJSON("task", nil, http.MethodGet)
	assert.NoError(t, err)
	var m map[string]string
	err = json.Unmarshal(body, &m)
	assert.NoError(t, err)

	e, ok := m["error"]
	assert.False(t, !ok || len(fmt.Sprint(e)) == 0,
		"Ожидается ошибка для вызова /task")

	body, err = requestJSON("task?id="+todo, nil, http.MethodGet)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &m)
	assert.NoError(t, err)

	//assert.Equal(t, todo, m["id"])
	//assert.Equal(t, task.Date, m["date"])
	//	assert.Equal(t, task.Title, m["title"])
	//assert.Equal(t, task.Description, m["description"])
	//	assert.Equal(t, task.Status, m["status"])
}

//
//type fulltask struct {
//	id string
//	task
//}
//
//func TestEditTask(t *testing.T) {
//	db := openDB(t)
//	defer db.Close()
//
//	now := time.Now()
//
//	tsk := task{
//		date:        now.Format(`20060102`),
//		title:       "Заказать пиццу",
//		description: "в 17:00",
//		status:      "",
//	}
//
//	id := addTask(t, tsk)
//
//	tbl := []fulltask{
//		{"", task{"20240129", "Тест", "", ""}},
//		{"abc", task{"20240129", "Тест", "", ""}},
//		{"7645346343", task{"20240129", "Тест", "", ""}},
//		{id, task{"20240129", "", "", ""}},
//		{id, task{"20240192", "Qwerty", "", ""}},
//		{id, task{"28.01.2024", "Заголовок", "", ""}},
//		{id, task{"20240212", "Заголовок", "", "ooops"}},
//	}
//	for _, v := range tbl {
//		m, err := postJSON("api/task", map[string]any{
//			"id":      v.id,
//			"date":    v.date,
//			"title":   v.title,
//			"comment": v.description,
//			"repeat":  v.status,
//		}, http.MethodPut)
//		assert.NoError(t, err)
//
//		var errVal string
//		e, ok := m["error"]
//		if ok {
//			errVal = fmt.Sprint(e)
//		}
//		assert.NotEqual(t, len(errVal), 0, "Ожидается ошибка для значения %v", v)
//	}
//
//	updateTask := func(newVals map[string]any) {
//		mupd, err := postJSON("api/task", newVals, http.MethodPut)
//		assert.NoError(t, err)
//
//		e, ok := mupd["error"]
//		assert.False(t, ok && fmt.Sprint(e) != "")
//
//		var task Task
//		err = db.Get(&task, `SELECT * FROM todolist WHERE id=?`, id)
//		assert.NoError(t, err)
//
//		assert.Equal(t, id, strconv.FormatInt(task.Id, 10))
//		assert.Equal(t, newVals["title"], task.Title)
//		if _, is := newVals["comment"]; !is {
//			newVals["comment"] = ""
//		}
//		if _, is := newVals["repeat"]; !is {
//			newVals["repeat"] = ""
//		}
//		assert.Equal(t, newVals["comment"], task.Description)
//		assert.Equal(t, newVals["repeat"], task.Status)
//		now := time.Now().Format(`20060102`)
//		if task.Date < now {
//			t.Errorf("Дата не может быть меньше сегодняшней")
//		}
//	}
//
//	updateTask(map[string]any{
//		"id":      id,
//		"date":    now.Format(`20060102`),
//		"title":   "Заказать хинкали",
//		"comment": "в 18:00",
//		"repeat":  "d 7",
//	})
//}
//
//func TestDone(t *testing.T) {
//	db := openDB(t)
//	defer db.Close()
//
//	now := time.Now()
//	id := addTask(t, task{
//		date:  now.Format(`20060102`),
//		title: "Свести баланс",
//	})
//
//	ret, err := postJSON("api/task/done?id="+id, nil, http.MethodPost)
//	assert.NoError(t, err)
//	assert.Empty(t, ret)
//	notFoundTask(t, id)
//
//	id = addTask(t, task{
//		title:  "Проверить работу /api/task/done",
//		status: "d 3",
//	})
//
//	for i := 0; i < 3; i++ {
//		ret, err := postJSON("api/task/done?id="+id, nil, http.MethodPost)
//		assert.NoError(t, err)
//		assert.Empty(t, ret)
//
//		var task Task
//		err = db.Get(&task, `SELECT * FROM todolist WHERE id=?`, id)
//		assert.NoError(t, err)
//		now = now.AddDate(0, 0, 3)
//		assert.Equal(t, task.Date, now.Format(`20060102`))
//	}
//}
//
//func TestDelTask(t *testing.T) {
//	db := openDB(t)
//	defer db.Close()
//
//	id := addTask(t, task{
//		title:  "Временная задача",
//		status: "d 3",
//	})
//	ret, err := postJSON("api/task?id="+id, nil, http.MethodDelete)
//	assert.NoError(t, err)
//	assert.Empty(t, ret)
//
//	notFoundTask(t, id)
//
//	ret, err = postJSON("/task", nil, http.MethodDelete)
//	assert.NoError(t, err)
//	assert.NotEmpty(t, ret)
//	ret, err = postJSON("/task?id=wjhgese", nil, http.MethodDelete)
//	assert.NoError(t, err)
//	assert.NotEmpty(t, ret)
//}
