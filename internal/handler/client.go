package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Rq4n/gofollow/internal/service"
	"github.com/Rq4n/gofollow/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ClientHandler struct {
	clientService service.ClientService
}

func NewClientHandler(clientService service.ClientService) *ClientHandler {
	return &ClientHandler{
		clientService: clientService,
	}
}

type ClientPayload struct {
	Name        string `json:"name"         validate:"required"`
	Email       string `json:"email"        validate:"required,email"`
	InvoiceLink string `json:"invoice_link" validate:"required, url"`
	SendDate    string `json:"send_date"    validate:"required"`
}

func (h *ClientHandler) HandleCreateClient(w http.ResponseWriter, r *http.Request) {
	var client ClientPayload
	if err := utils.ParseJSON(r, &client); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(client); err != nil {
		validationErrs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", validationErrs))
		return
	}

	sendDate, err := time.Parse("2006-01-02", client.SendDate)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid date format, expected YYYY-MM-DD"))
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, ErrUnauthorized)
		return
	}

	if err := h.clientService.CreateNewClient(r.Context(), userID, client.Name, client.Email, client.InvoiceLink, sendDate); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
