package main

import (
    "log"
    "net/http"
    "p2_GO_crud/handlers"
    "p2_GO_crud/middleware"
)

func main() {
    mux := http.NewServeMux()

    mux.Handle("/todos", http.HandlerFunc(handlers.TodosHandler))
    mux.Handle("/todos/", http.HandlerFunc(handlers.TodoByIDHandler))
    
    wrapped := middleware.Logging(mux)

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", wrapped))
}
