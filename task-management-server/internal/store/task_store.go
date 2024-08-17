package store

import (
	"fmt"
	"sync"
	"time"

	"example.com/task-management-server/internal/model"
	// "example.com/task-management-server/internal/helper"
)

type TaskStore struct {
	sync.Mutex
	tasks  map[int]model.Task
	nextId int
}

func New() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]model.Task),
		nextId: 0,
	}
}

func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) (int, error) {
	ts.Lock()
	defer ts.Unlock()

	if due.Before(time.Now()) {
		return 0, fmt.Errorf("due date must be later than the current time")
	}

	task := model.Task{
		Id:   ts.nextId,
		Text: text,
		Tags: tags,
		Due:  due,
	}

	ts.tasks[ts.nextId] = task
	ts.nextId++

	return task.Id, nil
}

func (ts *TaskStore) GetTask(id int) (model.Task, error) {
	ts.Lock()
	defer ts.Unlock()

	task, exists := ts.tasks[id]
	if !exists {
		return model.Task{}, fmt.Errorf("task with id = %d not found", id)
	}

	return task, nil
}

func (ts *TaskStore) GetAllTasks() []model.Task {
	ts.Lock()
	defer ts.Unlock()

	allTasks := make([]model.Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		allTasks = append(allTasks, task)
	}

	return allTasks

}

func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()

	_, exists := ts.tasks[id]
	if !exists {
		return fmt.Errorf("task with id = %d not found", id)
	}
	delete(ts.tasks, id)

	return nil
}

func (ts *TaskStore) DeletAllTasks() {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks = make(map[int]model.Task)
}

// helper module import not working properly
func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (ts *TaskStore) UpdateTask(id int, text *string, tags *[]string, due *time.Time) (model.Task, error) {

	ts.Lock()
	defer ts.Unlock()

	tasks, exists := ts.tasks[id]
	if !exists {
		return model.Task{}, fmt.Errorf("task with id = %d not found", id)
	}

	if text != nil {
		tasks.Text = *text
	}

	if tags != nil {
		for _, tag := range *tags {
			if !Contains(tasks.Tags, tag) {
				tasks.Tags = append(tasks.Tags, tag)
			}
		}
	}

	if due != nil {
		tasks.Due = *due
	}

	ts.tasks[id] = tasks

	return tasks, nil

}

func (ts *TaskStore) GetTasksByDueDate(year int, month time.Month, day int) []model.Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []model.Task
	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()
		if y == year && m == month && d == day {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

func (ts *TaskStore) GetTasksByTag(tag string) []model.Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []model.Task
	for _, task := range ts.tasks {
		for _, taskTag := range task.Tags {
			if taskTag == tag {
				tasks = append(tasks, task)
				break
			}
		}
	}

	return tasks
}
