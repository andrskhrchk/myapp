package httpserver

import "net/http"

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/sign-up", h.SignUp)
	mux.HandleFunc("POST /auth/sign-in", h.SignIn)
	return mux
}
