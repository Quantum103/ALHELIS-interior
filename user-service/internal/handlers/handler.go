package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"user-service/internal/middleware"
	"user-service/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) ShowLK(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())

	if !ok {
		http.Error(
			w,
			`{"error":"Пользователь не авторизован"}`,
			http.StatusUnauthorized,
		)
		return
	}

	log.Printf("PROFILE USER ID: %d", userID)

	profile, err := h.service.GetProfile(r.Context(), userID)

	log.Printf("PROFILE: %+v", profile)
	log.Printf("PROFILE ERROR: %v", err)

	if err != nil {
		http.Error(
			w,
			`{"error":"Профиль не найден"}`,
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(profile); err != nil {
		http.Error(
			w,
			`{"error":"Ошибка формирования ответа"}`,
			http.StatusInternalServerError,
		)
	}
}
func (h *UserHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"Пользователь не авторизован"}`, http.StatusUnauthorized)
		return
	}

	if err := h.service.CreateProfile(r.Context(), userID); err != nil {
		log.Printf("Ошибка создания профиля для userID=%d: %v", userID, err)
		http.Error(w, `{"error":"Ошибка создания профиля"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Профиль создан",
	})
}
