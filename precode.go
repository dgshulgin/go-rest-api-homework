package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...

// listTasks возвращает все задачи, содержащиеся в мапе tasks.
// При успешном выполнении возвращает 200 OK
//   - Task list JSONed
//   - OUT OF BR Empty task list
//
// В случае ошибки возвращает 500 Internale Server Error
//   - Unable to marshal map into JSON
func listTasks(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		debug.Printf("listTasks: Unable to marshal map into JSON\n")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// createTask создает новую задачу.
// Использует данные в теле запроса.
// Вновь созданная задача сохраняется в мапе tasks.
// При успешном выполнении возвращает 201 Created
//   - OUT OF BR Another task with the same description exists
//
// В случае ошибки возвращает 400 Bad Request
//   - OUT OF BR Both ID and Description are required with optional others
//   - Unable to get request's body
//   - Unable to parse invalid JSON
func createTask(w http.ResponseWriter, r *http.Request) {
	var reader bytes.Buffer
	_, err := reader.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		debug.Printf("createTask: Unable to get request's body\n")
		return
	}

	var task Task
	if err = json.Unmarshal(reader.Bytes(), &task); err != nil {
		http.Error(w, http.StatusText(400), 400)
		debug.Printf("createTask: Unable to parse invalid JSON\n")
		return
	}

	if len(task.ID) == 0 || len(task.Description) == 0 {
		http.Error(w, http.StatusText(400), 400)
		debug.Printf("createTask: Both ID and Description are required with optional others. Task creation rejected.\n")
		return
	}

	_, ok := tasks[task.ID]
	if ok { // Special behaviour left in business requirements
		debug.Printf("createTask: Another task with the same id exists. Overwrite.\n")
	}

	// now, set new task
	tasks[task.ID] = task
	debug.Printf("createTask: New task created %q.\n", task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// getTaskById возвращает задачу с указанным ID
// ID предоставляется как часть запроса.
// ID является ключем для задачи в мапе tasks.
// При успешном выполнении возвращает 200 OK
//   - Found and JSONed
//
// В случае ошибки возвращает 400 Bad Request
//   - Unable to find task with given ID
//   - Unable to marshall map's item into JSON
func getTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, http.StatusText(400), 400)
		debug.Printf("getTaskById: Unable to find task with given ID=%s\n", id)
		return
	}

	bytes, err := json.Marshal(task)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		debug.Printf("getTaskById: Unable to marshall map's item into JSON\n")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// deleteTaskById удаляет задачу с указанным ID из хранилища (мап tasks)
// ID предоставляется как часть запроса.
// ID является ключем для задачи в мапе tasks.
// При успешном выполнении возвращает 200 OK
//   - Deletion for ID succeded
//
// В случае ошибки возвращает 400 Bad Request
//   - Unable to find task with given ID
func deleteTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, ok := tasks[id]
	if !ok {
		http.Error(w, http.StatusText(400), 400)
		debug.Printf("deleteTaskById: Unable to find task with given ID=%s\n", id)
		return
	}

	// task deletion
	delete(tasks, id)
	debug.Printf("deleteTaskById: Task id=%s was deleted.\n", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

var debug *log.Logger

func main() {
	r := chi.NewRouter()

	debug = log.New(os.Stdout, "DEBUG\t", log.Ltime)
	// здесь регистрируйте ваши обработчики
	// ...
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", listTasks)   // GET /tasks
		r.Post("/", createTask) // POST /tasks

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", getTaskById)       // GET /tasks/{taskID}
			r.Delete("/", deleteTaskById) // DELETE /tasks/{taskID}
		})
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		debug.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
