package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"
    "sort"
)

type ByCreatedAt []Todo

func (b ByCreatedAt) Len() int {
    return len(b)
}

func (b ByCreatedAt) Less(i, j int) bool {
    return b[i].CreatedAt < b[j].CreatedAt
}

func (b ByCreatedAt) Swap(i, j int) {
    b[i], b[j] = b[j], b[i]
}

type Todo struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    CreatedAt string `json:"createdAt"`
    Done  bool   `json:"done"` 
}

var todos []Todo
var nextId = 3

func main() {
    todos = append(todos, Todo{1, "first", time.Now().Format("2006-01-02 15:04:05"), false})
    time.Sleep(1 * time.Second)
    todos = append(todos, Todo{2, "second", time.Now().Format("2006-01-02 15:04:05"), false})

    http.HandleFunc("/todos", withLogging(handleTodos))
    http.HandleFunc("/todos/", withLogging(handleTodoByID)) //note trailing slash
    http.HandleFunc("/health", withLogging(handleHealth)) //note trailing slash

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func withLogging(next http.HandlerFunc) http.HandlerFunc{
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next(w, r)
    }
}


func handleHealth(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    } else {
        response := map[string]string{"status" : "ok"}
        writeJSON(w, response)
    }
}

func handleTodos(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        
        done := r.URL.Query().Get("done")

        if done == "true" {
            var trueTodos []Todo
            for _, t := range todos {
                if t.Done == true {
                    trueTodos = append(trueTodos, t)
                }
            }
            writeJSON(w, trueTodos)
        } else if done == "" {
            sort.Sort(ByCreatedAt(todos))
            writeJSON(w, todos)
        } else {
            http.Error(w, "done value should be true", http.StatusBadRequest)
        }

    case http.MethodPost:
        var t Todo
        if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
            http.Error(w, "invalid json", http.StatusBadRequest)
            return
        }
        if t.Title == "" {
            http.Error(w, "Title empty", http.StatusBadRequest)
            return
        }
        t.ID = nextId
        t.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
        nextId++
        todos = append(todos, t)
        w.WriteHeader(http.StatusCreated)
        writeJSON(w, t)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

func handleTodoByID(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
    id, err := strconv.Atoi(idStr)

    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest) 
        return
    }

    index := -1
    for i, t := range todos {
        if t.ID == id {
            index = i
            break;
        }
    }
    if index == -1 {
        http.Error(w, "not found", http.StatusBadRequest)
        return
    }

    switch r.Method {
    case http.MethodGet:
        writeJSON(w, todos[index])
    case http.MethodPut:
        var updated Todo
        if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
            http.Error(w, "invalid json", http.StatusBadRequest)
            return
        }
        updated.ID = id
        todos[index] = updated
        writeJSON(w, updated)
    case http.MethodDelete:
        todos = append(todos[:index], todos[index+1:]...)
        w.WriteHeader(http.StatusNoContent)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

func writeJSON(w http.ResponseWriter, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(v)
}

