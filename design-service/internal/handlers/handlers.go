package handlers

import (
	"fmt"
	"net/http"
)

func DashboardFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
}
