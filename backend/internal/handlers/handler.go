package handlers

import (
	"auth/backend/internal/services"
	"auth/backend/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	s *services.Service
}

func NewHandler(s *services.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", h.LoginHandler).Methods("POST")
	r.HandleFunc("/refresh", h.RefreshHandler).Methods("POST")

	return r
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOGIN")
	id := utils.GetUserID(r)
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
	}

	ipAddr, err := utils.GetIP(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	loginRes, err := h.s.Login(id, ipAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(loginRes)

}

func (h *Handler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]string

	userID := utils.GetUserID(r)
	if userID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	refresh, ok := body["refresh"]
	if !ok {
		http.Error(w, "refresh is required", http.StatusBadRequest)
		return
	}

	ipAddr, _ := utils.GetIP(r)

	res, err := h.s.Refresh(refresh, userID, ipAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(res)

}
