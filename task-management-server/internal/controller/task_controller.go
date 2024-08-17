package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"example.com/task-management-server/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type TaskController struct {
	service  *service.TaskService
	validate *validator.Validate
}

func NewTaskController(service *service.TaskService) *TaskController {
	return &TaskController{
		service:  service,
		validate: validator.New(),
	}
}

type RequestTask struct {
	Text string    `json:"text" validate:"required,min=1"`
	Tags []string  `json:"tags" validate:"dive,required"`
	Due  time.Time `json:"due" validate:"required"`
}

func (tc *TaskController) CreateTaskHandler(w http.ResponseWriter, req *http.Request) {
	var rt RequestTask

	if err := json.NewDecoder(req.Body).Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := tc.validate.Struct(rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if rt.Due.Before(time.Now()) {
		http.Error(w, "due date must be later than the current time", http.StatusBadRequest)
		return
	}

	id, err := tc.service.CreateTask(rt.Text, rt.Tags, rt.Due)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (tc *TaskController) GetTaskHandler(w http.ResponseWriter, req *http.Request) {

	idStr := mux.Vars(req)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := tc.service.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)

	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (tc *TaskController) GetAllTasksHandler(w http.ResponseWriter, req *http.Request) {

	tasks := tc.service.GetAllTasks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (tc *TaskController) DeleteTaskHandler(w http.ResponseWriter, req *http.Request) {

	idStr := mux.Vars(req)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := tc.service.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)

	}
	w.WriteHeader(http.StatusNoContent)
}

func (tc *TaskController) DeleteAllTasksHandler(w http.ResponseWriter, req *http.Request) {

	tc.service.DeletAllTasks()
	w.WriteHeader(http.StatusNoContent)
}

func (tc *TaskController) GetTaskByTagHandler(w http.ResponseWriter, req *http.Request) {

	tag := mux.Vars(req)["tag"]
	tasks := tc.service.GetTasksByTag(tag)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (tc *TaskController) GetTaskByDueDateHandler(w http.ResponseWriter, req *http.Request) {

	yearStr := req.URL.Query().Get("year")
	monthStr := req.URL.Query().Get("month")
	dayStr := req.URL.Query().Get("day")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "invalid year", http.StatusBadRequest)
		return
	}

	monthInt, err := strconv.Atoi(monthStr)
	if err != nil {
		http.Error(w, "invalid month", http.StatusBadRequest)
		return
	}
	month := time.Month(monthInt)

	day, err := strconv.Atoi(dayStr)
	if err != nil {
		http.Error(w, "invalid day", http.StatusBadRequest)
		return
	}

	tasks := tc.service.GetTasksByDueDate(year, month, day)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (tc *TaskController) UpdateTaskHandler(w http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var taskRequest struct {
		Text *string    `json:"text,omitempty"`
		Tags *[]string  `json:"tags,omitempty"`
		Due  *time.Time `json:"due,omitempty"`
	}

	if err := json.NewDecoder(req.Body).Decode(&taskRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTask, err := tc.service.UpdateTask(id, taskRequest.Text, taskRequest.Tags, taskRequest.Due)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)

}
