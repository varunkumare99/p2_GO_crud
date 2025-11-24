package store

import (
    "p2_GO_crud/models"
)

var todos []models.Todo
var nextID = 1

func GetAllTodos() []models.Todo {
    return todos
}

func GetTodoByID(id int) (models.Todo, bool) {
    for _, t := range todos {
        if t.ID == id {
            return t, true
        }
    }
    return models.Todo{}, false
}

func CreateTodo(title string) models.Todo {
    todo := models.Todo {
        ID: nextID,
        Title: title,
        Done: false,
    }

    nextID++
    todos = append(todos, todo)
    return todo
}


func UpdateTodo(id int, updated models.Todo) bool {
    for i := range todos {
        if todos[i].ID == id {
            updated.ID = id
            todos[i] = updated
            return true
        }
    }
    return false
}

func DeleteTodo(id int) bool {
    for i := range todos {
        if todos[i].ID == id {
            todos = append(todos[:i], todos[i+1:]...)
            return true
        }
    }
    return false
}
