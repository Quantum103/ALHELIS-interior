package main

import (
	"design-service/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/dashboard", handlers.DashboardFunc)

	http.ListenAndServe(":8082", r)
}
