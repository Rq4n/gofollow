package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ClientService interface {
	CreateNewClient(ctx context.Context, userID uuid.UUID, name, email, invoice_link string, sendDate time.Time) error
}

type ClientHandler struct {
	clientService ClientService
}

func NewClientHandler(clientService ClientService) *ClientHandler {
	return &ClientHandler{
		clientService: clientService,
	}
}

func (h *ClientHandler) HandleCreateClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	var req Client
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	sendDate, err := time.Parse("2006-01-02", req.SendDate)
	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, ErrNotAuthorized.Error(), http.StatusUnauthorized)
		return
	}

	if err := h.clientService.CreateNewClient(r.Context(), userID, req.Name, req.Email, req.InvoiceLink, sendDate); err != nil {
		http.Error(w, ErrInternalServerErr.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "client created successfully",
		"name":    req.Name,
		"email":   req.Email,
		"invoice": req.InvoiceLink,
	})
}
