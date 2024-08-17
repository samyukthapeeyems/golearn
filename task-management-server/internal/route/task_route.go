package route

import (
	"example.com/task-management-server/internal/controller"
	"example.com/task-management-server/internal/service"
	"example.com/task-management-server/internal/store"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	taskStore := store.New()
	taskService := service.NewTaskService(taskStore)
	taskController := controller.NewTaskController(taskService)

	router := mux.NewRouter()

	router.HandleFunc("/tasks", taskController.CreateTaskHandler).Methods("POST")
	router.HandleFunc("/tasks", taskController.GetAllTasksHandler).Methods("GET")
	router.HandleFunc("/tasks/{id:[0-9]+}", taskController.GetTaskHandler).Methods("GET")
	router.HandleFunc("/tasks/{id:[0-9]+}", taskController.DeleteTaskHandler).Methods("DELETE")
	router.HandleFunc("/tasks", taskController.DeleteAllTasksHandler).Methods("DELETE")
	router.HandleFunc("/tasks/{id:[0-9]+}", taskController.UpdateTaskHandler).Methods("PATCH")
	router.HandleFunc("/tag/{tag}", taskController.GetTaskByTagHandler).Methods("GET")
	router.HandleFunc("/due", taskController.GetTaskByDueDateHandler).Methods("GET")

	return router
}
