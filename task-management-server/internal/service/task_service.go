package service

import (
	"fmt"
	"time"

	"example.com/task-management-server/internal/model"
	"example.com/task-management-server/internal/store"
)

type TaskService struct {
	store *store.TaskStore
}

func NewTaskService(store *store.TaskStore) *TaskService {
	return &TaskService{store: store}
}

func (ts *TaskService) CreateTask(text string, tags []string, due time.Time) (int, error) {
	task := &model.Task{
		Text: text,
		Tags: tags,
		Due:  due,
	}

	if err := task.Validate(); err != nil {
		return 0, err
	}

	if due.Before(time.Now()) {
		return 0, fmt.Errorf("due date must be later than the current time")
	}

	return ts.store.CreateTask(text, tags, due)
}

func (ts *TaskService) UpdateTask(id int, text *string, tags *[]string, due *time.Time) (model.Task, error) {
	existingTask, err := ts.store.GetTask(id)
	if err != nil {
		return model.Task{}, err
	}

	if text != nil {
		existingTask.Text = *text
	}
	if tags != nil {
		existingTask.Tags = *tags
	}
	if due != nil {
		if due.Before(time.Now()) {
			return model.Task{}, fmt.Errorf("due date must be later than the current time")
		}
		existingTask.Due = *due
	}

	if err := existingTask.Validate(); err != nil {
		return model.Task{}, err
	}

	return ts.store.UpdateTask(id, text, tags, due)
}

func (ts *TaskService) GetTask(id int) (model.Task, error) {
	return ts.store.GetTask(id)
}

func (ts *TaskService) GetAllTasks() []model.Task {
	return ts.store.GetAllTasks()
}

func (ts *TaskService) DeleteTask(id int) error {
	return ts.store.DeleteTask(id)
}

func (ts *TaskService) DeletAllTasks() {
	ts.store.DeletAllTasks()
}

func (ts *TaskService) GetTasksByDueDate(year int, month time.Month, day int) []model.Task {
	return ts.store.GetTasksByDueDate(year, month, day)
}

func (ts *TaskService) GetTasksByTag(tag string) []model.Task {
	return ts.store.GetTasksByTag(tag)
}
