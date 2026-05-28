package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/config"
	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/middleware"
	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/users"
	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/validator"
)

type AuthHandler struct {
	userService *users.Service
	config      *config.Config
}

func NewAuthHandler(userService *users.Service, cfg *config.Config) *AuthHandler {
	return &AuthHandler{userService: userService, config: cfg}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req validator.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.Register(r.Context(), req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"data":    user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req validator.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, user, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"data": map[string]interface{}{
			"token": token,
			"user":  user,
		},
	})
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := int64(middleware.GetUserID(r.Context()))
	
	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": user,
	})
}
