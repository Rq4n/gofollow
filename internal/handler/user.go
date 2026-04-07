package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Rq4n/gofollow/internal/auth"
	"github.com/Rq4n/gofollow/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	var req UserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userService.CreateNewUser(r.Context(), req.Email, req.Password); err != nil {
		http.Error(w, ErrInternalServerErr.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "user created successfully",
	})
}

func (h *UserHandler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	var req UserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.userService.GetUserByName(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, ErrInternalServerErr.Error(), http.StatusInternalServerError)
		return
	}

	tokenStr, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, ErrInternalServerErr.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{
		"token": tokenStr,
	})
}
