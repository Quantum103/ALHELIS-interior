package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	userIDStr := r.Header.Get("X-User-ID")

	if userIDStr == "" {
		http.Error(w, "user id not provided", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	profile, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(profile); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
