package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "p2_GO_crud/models"
    "p2_GO_crud/store"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        w.WriteHeader(http.StatusOK)
        writeJSON(w, "boss I'm in good health")
    } else {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

func TodosHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        writeJSON(w, store.GetAllTodos())

    case http.MethodPost:
        var data struct {
            Title string `json:"title"`
        }

        if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
            http.Error(w, "invalid json", http.StatusBadRequest)
            return
        }

        if strings.TrimSpace(data.Title) == "" {
            http.Error(w, "title cannot be empty", http.StatusBadRequest)
            return
        }

        newTodo := store.CreateTodo(data.Title)

        w.Header().Set("Location", "/todos/"+strconv.Itoa(newTodo.ID))
        w.WriteHeader(http.StatusCreated)
        writeJSON(w, newTodo)

    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

func TodoByIDHandler(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
    id, err := strconv.Atoi(idStr)

    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }

    switch r.Method {
    case http.MethodGet:
        todo, ok := store.GetTodoByID(id)
        if !ok {
            http.Error(w, "not found", http.StatusNotFound)
            return
        }
        writeJSON(w, todo)

    case http.MethodPut:
        var updated  models.Todo
        if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
            http.Error(w, "invalid json", http.StatusBadRequest)
            return
        }

        if !store.UpdateTodo(id, updated) {
            http.Error(w, "not found", http.StatusNotFound)
            return
        }

        writeJSON(w, updated)
        
    case http.MethodDelete:
        if !store.DeleteTodo(id) {
            http.Error(w, "not found", http.StatusNotFound)
            return
        }
        w.WriteHeader(http.StatusNoContent)

    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func writeJSON(w http.ResponseWriter, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(v)
}
            
