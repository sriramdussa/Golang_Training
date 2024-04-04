package main

import ("fmt"
"net/http")

import "github.com/gorilla/mux"

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/", handler).Methods("GET")

    router.Use(loggingMiddleware)

    fmt.Println("Server is running on port 8080...")
    http.ListenAndServe(":8080", router)
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Request: %s %s\n", r.Method, r.RequestURI)
        next.ServeHTTP(w, r)
    })
}
