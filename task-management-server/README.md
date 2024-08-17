## Task Management Server

### Introduction

This project is a RESTful API server for managing tasks. It includes features such as task creation, retrieval, updating, and deletion, all implemented using Go. The server also supports filtering tasks by due date and tags.In a realistic application, the TaskStore here would likely be an interface that several backends can implement, but for our simple example the current API is sufficient. If you'd like, you could implement a TaskStore using something like MongoDB or MySQL.



### Features

- Create a new task
- Get task by ID
- Update an existing task (partially)
- Delete a task
- Retrieve all tasks
- Filter tasks by due date
- Filter tasks by tag

### Installation

To run this project, you need to have Go installed. Follow these steps:

1. Clone the repository:

    ```bash
    git clone https://github.com/samyukthapeeyems/task-management-server.git
    ```

2. Navigate to the project directory:

    ```bash
    cd task-management-server
    ```

3. Build and run the server:

    ```bash
    go run cmd/main.go
    ```

4. Test the API using tools like Postman or Hoppscotch.

### Usage

- **Create Task**: `POST /tasks`
- **Get Task**: `GET /tasks/{id}`
- **Update Task**: `PATCH /tasks/{id}`
- **Delete Task**: `DELETE /tasks/{id}`
- **Get All Tasks**: `GET /tasks`
- **Get Tasks by Due Date**: `GET /tasks/due?year={year}&month={month}&day={day}`
- **Get Tasks by Tag**: `GET /tasks/tag/{tag}`
