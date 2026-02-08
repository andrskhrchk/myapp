package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	jwt "github.com/andrskhrchk/myapp/pkg/jwt"

	"github.com/andrskhrchk/myapp/internal/services/auth"
	"github.com/andrskhrchk/myapp/internal/transport/dto"
)

type Handler struct {
	services     *auth.AuthService
	tokenManager *jwt.TokenManager
}

func NewHandler(services *auth.AuthService, tokenManager *jwt.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var input dto.RegisterDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	log.Println(input.Email, input.Password)

	if input.Email == "" || input.Password == "" {
		http.Error(w, "email and password is required", http.StatusBadRequest)
		return
	}

	user, token, err := h.services.Register(r.Context(), &input)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to register user %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var loginData dto.LoginDTO

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if loginData.Email == "" || loginData.Password == "" {
		http.Error(w, "email and password is required", http.StatusBadRequest)
		return
	}

	user, token, err := h.services.Login(r.Context(), &loginData)

	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user":  user,
	})
}
