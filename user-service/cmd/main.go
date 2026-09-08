package main

import (
	"log"
	"net/http"

	"user-service/internal/database"
	"user-service/internal/handlers"
	"user-service/internal/repo"
	"user-service/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	cfg := database.LoadConfigFromEnv()

	db, err := database.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repo.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	r := mux.NewRouter()

	r.HandleFunc(
		"/api/profile",
		userHandler.ShowLK,
	).Methods(http.MethodGet)

	log.Println("User Service запущен на :8082")

	if err := http.ListenAndServe(":8082", r); err != nil {
		log.Fatal(err)
	}
}

/*
маршруты для user-service

const API = {
    profile: "/api/profile",
    projects: "/api/profile/projects",
    favorites: "/api/profile/favorites",
    updateProfile: "/api/profile",
    changePassword: "/api/auth/password",
    logout: "/api/auth/logout"
};
*/
