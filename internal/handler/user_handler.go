package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Rq4n/gofollow/internal/auth"
	"github.com/Rq4n/gofollow/internal/repository"
)

type UserService interface {
	CreateNewUser(ctx context.Context, email, password string) error
	GetUserByName(ctx context.Context, email, password string) (*repository.GetUserByNameRow, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	var req User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userService.CreateNewUser(r.Context(), req.Name, req.Password); err != nil {
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

	var req User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.userService.GetUserByName(r.Context(), req.Name, req.Password)
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
