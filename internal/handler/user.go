package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Rq4n/gofollow/internal/auth"
	"github.com/Rq4n/gofollow/internal/service"
	"github.com/Rq4n/gofollow/internal/utils"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		validationErrs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", validationErrs))
		return
	}

	if err := h.userService.CreateNewUser(r.Context(), user.Email, user.Password); err != nil {
		switch {
		case errors.Is(err, service.ErrUserAlreadyExists):
			utils.WriteError(w, http.StatusConflict, err)
		default:
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *UserHandler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	var user LoginUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		validationErrs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", validationErrs))
		return
	}

	resp, err := h.userService.GetUserByEmail(r.Context(), user.Email, user.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials),
			errors.Is(err, service.ErrUserNotFound):
			utils.WriteError(w, http.StatusUnauthorized, err)
		default:
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	tokenStr, err := auth.GenerateToken(resp.ID, user.Email)
	if err != nil {
		http.Error(w, ErrInternalServerErr.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{"token": tokenStr})
}
